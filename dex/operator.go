package dex

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/Dvisacker/matching-engine/dex/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Operator manages the transaction queue that will eventually be
// sent to the exchange contract. The Operator Wallet must be equal to the
// account that initially deployed the exchange contract or an address with operator rights
// on the contract
// - Admin is the wallet that sends the transactions to the exchange smart-contract
// - EthereumClient contains
// - Params contains the
// - Chain seems to be deprecated
// - TxLogs seems to be deprecated
// - ErrorChannel, TradeChannel and CancelOrderChannel listen and pipes smart-contract to the engine handler goroutine
// - ErrorLogs, TradeLogs and CancelOrderLogs listen on the
// - OrderTradePairs contains a mapping of the current trade hashes that have been sent to the contract.
// - ErrorLogs, TradeLogs and CancelOrderLogs don't seem to be used as of now
type Operator struct {
	Admin              *Wallet
	Exchange           *Exchange
	Engine             *TradingEngine
	EthereumClient     *ethclient.Client
	Params             *OperatorParams
	Chain              bind.ContractBackend
	TxQueue            []Hash
	ErrorChannel       chan *interfaces.ExchangeLogError
	TradeChannel       chan *interfaces.ExchangeLogTrade
	CancelOrderChannel chan *interfaces.ExchangeLogCancelOrder
	OrderTradePairs    map[Hash]*OrderTradePair
}

// NewOperator creates a new operator struct. It creates an exchange contract instance from the
// provided address. The error and trade events are received in the ErrorChannel and TradeChannel.
// Upon receiving errors and trades in their respective channels, event payloads are sent to the
// associated order maker and taker sockets through the through the event channel on the Order and Trade struct.
// In addition, an error event cancels the trade in the trading engine and makes the order available again.
func NewOperator(config *OperatorConfig) (*Operator, error) {
	op := &Operator{}

	rpcClient, err := rpc.DialWebsocket(context.Background(), config.OperatorParams.rpcURL, "")
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)
	ex, err := NewExchange(config.Admin, config.Exchange, client)
	if err != nil {
		return nil, err
	}

	op.Admin = config.Admin
	op.Params = config.OperatorParams
	op.Exchange = ex
	op.EthereumClient = client
	op.OrderTradePairs = make(map[Hash]*OrderTradePair)
	op.TxQueue = make([]Hash, 0)

	op.ErrorChannel, err = op.Exchange.ListenToErrorEvents()
	if err != nil {
		return nil, err
	}

	op.TradeChannel, err = op.Exchange.ListenToTrades()
	if err != nil {
		return nil, err
	}

	err = op.Validate()
	if err != nil {
		return nil, err
	}

	// Bug: In certain cases, the trade channel seems to be receiving additional unexpected trades.
	// In the case TestSocketExecuteOrder (in file socket_test.go) is run on its own, everything is working correctly.
	// However, in the case TestSocketExecuteOrder is run among other tests, some tradeLogs do not correspond to an
	// order hash in the ordertrade mapping. I suspect this is because the event listener catches events from previous
	// tests. It might be helpful to see how to listen to events from up to a certain block.
	go func() {
		for {
			select {
			case errLog := <-op.ErrorChannel:
				tradeHash := errLog.TradeHash
				otp, ok := op.OrderTradePairs[tradeHash]
				if !ok {
					log.Printf("Could not retrieve initial hash")
					return
				}

				t := otp.trade
				o := otp.order
				errID := errLog.ErrorId

				if otp.order.events != nil {
					otp.order.events <- otp.order.NewOrderTxError(t, errID)
				}

				if otp.trade.events != nil {
					otp.trade.events <- otp.trade.NewTradeTxError(o, errID)
				}

				if op.Engine != nil {
					if err := op.Engine.CancelTrade(t); err != nil {
						log.Printf("Could not update order")
						return
					}
				}

			case tradeLog := <-op.TradeChannel:
				otp, ok := op.OrderTradePairs[tradeLog.TradeHash]
				if !ok {
					log.Printf("Could not retrieve initial hash")
					return
				}

				// only execute the next transaction in the queue when this transaction is mined
				go func() {
					_, err := op.WaitMined(otp.tx)
					if err != nil {
						log.Printf("Could not execute trade: %v\n", err)
					}

					if otp.order.events != nil {
						otp.order.events <- otp.order.NewOrderTxSuccess(otp.trade, otp.tx)
					}

					if otp.trade.events != nil {
						otp.trade.events <- otp.trade.NewTradeTxSuccess(otp.order, otp.tx)
					}

					op.TxQueue = op.TxQueue[1:]
					if len(op.TxQueue) > 0 {
						newOtp, ok := op.OrderTradePairs[op.TxQueue[0]]
						if !ok {
							log.Printf("Error parsing queue")
						}
						_, err := op.ExecuteTrade(newOtp.order, newOtp.trade)
						if err != nil {
							log.Printf("Could not execute trade: %v\n", err)
						}
					}
				}()

			}
		}
	}()

	return op, nil
}

// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
func (op *Operator) AddTradeToExecutionList(o *Order, t *Trade) error {
	err := t.Validate()
	if err != nil {
		return err
	}

	op.OrderTradePairs[t.Hash] = &OrderTradePair{order: o, trade: t, tx: nil}
	op.TxQueue = append(op.TxQueue, t.Hash)

	if len(op.TxQueue) == 1 {
		go op.ExecuteTrade(o, t)
	}

	if len(op.TxQueue) == 10 {
		return errors.New("Transaction queue is full")
	}

	return nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (op *Operator) ExecuteTrade(o *Order, t *Trade) (*types.Transaction, error) {
	tx, err := op.Exchange.Trade(o, t)
	if err != nil {
		return nil, err
	}

	op.OrderTradePairs[t.Hash].tx = tx

	if o.events != nil {
		o.events <- o.NewOrderExecutedEvent(tx)
	}

	if o.events != nil {
		t.events <- t.NewTradeExecutedEvent(tx)
	}

	return tx, nil
}

// WaitMined blocks the current the current goroutine execution until the given transaction gets mined on the blockchain
// the operator is connected to
func (op *Operator) WaitMined(tx *types.Transaction) (*types.Receipt, error) {
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, op.EthereumClient, tx)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// Validate checks that the operator configuration is sufficient.
func (op *Operator) Validate() error {
	balance, err := op.EthereumClient.PendingBalanceAt(context.Background(), op.Admin.Address)
	if err != nil {
		return err
	}

	if balance.Cmp(op.Params.minBalance) == -1 {
		return errors.New("Balance is insufficient")
	}

	return nil
}

// SetDefaultTxOptions resets the transaction value to 0
func (op *Operator) SetDefaultTxOptions() {
	op.Exchange.TxOptions.Value = big.NewInt(0)
}

// SetTxValue sets the transaction ether value
func (op *Operator) SetTxValue(value *big.Int) {
	op.Exchange.TxOptions.Value = value
}

// SetCustomSender updates the sender address address to the exchange contract
func (op *Operator) SetCustomSender(w *Wallet) {
	op.Exchange.TxOptions = bind.NewKeyedTransactor(w.PrivateKey)
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (op *Operator) SetFeeAccount(account Address) (*types.Transaction, error) {
	tx, err := op.Exchange.SetFeeAccount(account)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (op *Operator) SetOperator(account Address, isOperator bool) (*types.Transaction, error) {
	tx, err := op.Exchange.SetOperator(account, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetWithdrawalSecurityPeriod sets the period after which a non-operator address can send
// a transaction to the exchange smart-contract to withdraw their funds. This acts as security mechanism
// to prevent the operator of the exchange from holding funds
func (op *Operator) SetWithdrawalSecurityPeriod(p *big.Int) (*types.Transaction, error) {
	tx, err := op.Exchange.SetWithdrawalSecurityPeriod(p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositEther deposits ether into the exchange smart-contract. A priori this function is not supposed
// to be called by the exchange operator
func (op *Operator) DepositEther(val *big.Int) (*types.Transaction, error) {
	tx, err := op.Exchange.DepositEther(val)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositToken deposits tokens into the exchange smart-contract. A priori this function is not supposed
// to be called by the exchange operator
func (op *Operator) DepositToken(token Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := op.Exchange.DepositToken(token, amount)
	if err != nil {
		return nil, err
	}

	return tx, err
}

// TokenBalance returns the Exchange token balance of the given token at the given account address.
// Note: This is not the token BalanceOf() function, it's the balance of tokens that have been deposited
// in the exchange smart contract.
func (op *Operator) TokenBalance(account Address, token Address) (*big.Int, error) {
	b, err := op.Exchange.TokenBalance(account, token)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// EtherBalance returns the Exchange ether balance of the given account address.
// Note: This is not the current ether balance of the given ether address. It's the balance of ether
// that has been deposited in the exchange smart contract.
func (op *Operator) EtherBalance(account Address) (*big.Int, error) {
	b, err := op.Exchange.EtherBalance(account)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// WithdrawalSecurityPeriod is the period after which a non-operator account can withdraw their funds from
// the exchange smart contract.
func (op *Operator) WithdrawalSecurityPeriod() (*big.Int, error) {
	p, err := op.Exchange.WithdrawalSecurityPeriod()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (op *Operator) FeeAccount() (Address, error) {
	account, err := op.Exchange.FeeAccount()
	if err != nil {
		return Address{}, err
	}

	return account, nil
}

// Operator returns true if the given address is an operator of the exchange and returns false otherwise
func (op *Operator) Operator(addr Address) (bool, error) {
	isOperator, err := op.Exchange.Operator(addr)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

// SecurityWithdraw executes a security withdraw transaction. Security withdraw transactions can only be
// executed after the security withdrawal period has ended. A priori, this function should not be called
// by the operator account itself
func (op *Operator) SecurityWithdraw(w *Wallet, token Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := op.Exchange.SecurityWithdraw(w, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Withdraw executes a normal withdraw transaction. This withdraws tokens or ether from the exchange
// and returns them to the payload Receiver. Only an operator account can send a withdraw
// transaction
func (op *Operator) Withdraw(w *Withdrawal) (*types.Transaction, error) {
	tx, err := op.Exchange.Withdraw(w)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
