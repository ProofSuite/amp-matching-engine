package ethereum

import (
	"encoding/json"
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/dex/interfaces"
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
	WalletService *services.walletService,
	TxService *services.txService,
	TradeServoce *services.tradeService
	Exchange *contracts.Exchange,
}

type OperatorMessage struct {
	MessageType string,
	Order *types.Order,
	Trade *types.Trade
	ErrID int
}

type PendingTradeMessage struct {
	Order *types.Order,
	Trade *types.Trade
}

// NewOperator creates a new operator struct. It creates an exchange contract instance from the
// provided address. The error and trade events are received in the ErrorChannel and TradeChannel.
// Upon receiving errors and trades in their respective channels, event payloads are sent to the
// associated order maker and taker sockets through the through the event channel on the Order and Trade struct.
// In addition, an error event cancels the trade in the trading engine and makes the order available again.
func InitOperator(
	walletService *services.WalletService,
	txService *services.TxService,
	tradeService *services.TradeService
	exchange *contracts.Exchange,
) {
	op := &Operator {
		WalletService: walletService,
		TxService: txService,
		TradeService: tradeService,
		Exchange: *contracts
	}

	tradeEvents, err := exchange.ListenToTrades()
	if err != nil {
		return
	}

	errorEvents, err := exchange.ListenToErrors()
	if err != nil {
		return
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
			case event := <-errorEvents:
				tradeHash := event.TradeHash
					tr, ok := pendingTradeService.getTrade(tradeHash)
					if !ok {
						log.Printf("Could not retrieve hash")
						return
					}

					err := PublishTxErrorMessage(tr, errID)
					if err != nil {
						log.Printf("Could not publish tx error message")
					}

					err := PublishCancelTradeMessage(tr)
					if err != nil {
						log.Printf("Could not publish cancel trade message")
					}

			case event := <-tradeEvents:
				tr, ok := tradeService.getTradeByHash(event.TradeHash)
				if !ok {
					log.Printf("Could not retrieve initial hash")
					return
				}

				// only execute the next transaction in the queue when this transaction is mined
				go func() {
					_, err := op.EthereumService.WaitMined(tr.tx)
					if err != nil {
						log.Printf("Could not execute trade: %v\n", err)
					}

					err := op.PublishTradeSuccessMessage(tr)
					if err != nil {
						log.Printf("Could not publish order success message")
					}

					ch := getChannel("PENDING_TRADES")
					q := getQueue(ch, "PENDING_TRADES")

					length := ch.QueueInspect.Messages
					if length > 0 {
						msg, err := ch.Get(
							q.Name,
							"",
							true,
							false,
							false,
							false,
							nil
						)

						var pendingTrade PendingTradeMessage
						err := json.Unmarshal(msg.Body, &pendingTrade)
						if err != nil {
							log.Printf("Could not executed trade: %v\n", err)
						}

						_, err := op.ExecuteTrade(pendingTrade.Order, pendingTrade.Trade)
						if err != nil {
							log.Printf("Could not execute trade: %v", err)
						}
					}
				}()
			}
		}
	}()

	return op, nil
}

func (op *Operator) SubscribeOperatorMessages(fn func(*OperatorMessage) error) error {
	ch := getChannel("OPERATOR_SUB")
	q := getQueue("TX_MESSAGES")

	go func() {
		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil
		)

		if err != nil {
			log.Fataf("Failed to register a consumer: %s", err)
		}

		forever := make(chan bool)

		go func() {
			for m := msgs {
				var om *OperatorMessage
				err := json.Unmarshal(m.Body, &om)
				if err != nil {
					return err
					continue
				}
				go fn(om)
			}
		}()

		<-forever
	}()
	return nil
}

func (op *Operator) PublishTxErrorMessage(tr, errID) error {
	msg := &OperatorMessage{
		MessageType: "TX_ERROR_MESSAGE"
		Trade: tr,
		ErrID: errID
	}

	err := op.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

func (op *Operator) PublishTradeCancelMessage(tr) error {
	msg := &OperatorMessage{
		MessageType: "TRADE_CANCEL_MESSAGE",
		Trade: tr,
	}

	err := op.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

func (op *Operator) PublishTradeSuccessMessage(tr) error {
	msg := &OperatorMessage{
		MessageType: "TRADE_SUCCESS_MESSAGE",
		Trade: tr,
	}

	err := op.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

func (op *Operator) Publish(msg *OperatorMessage) error {
	ch := getChannel("OPERATOR_PUB")
	q := getQueue(ch, "TX_MESSAGES")

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		ampq.Publish{
			ContentType: "text/json",
			Body: bytes
		}
	)

	if err != nil {
		log.Printf("Failed to publish message %s: %s" msg.MessageType, err)
		return err
	}

	return nil
}


// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
func (op *Operator) QueueTrade(o *Order, t *Trade) error {
	err := t.Validate()
	if err != nil {
		return err
	}

	ch := getChannel("tradeTxs")
	q := getQueue(ch, "tradeTxs")

	bytes, err := json.Marshal(t)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body: bytes
		})

	length := ch.QueueInspect.Messages
	if length == 1 {
		op.ExecuteTrade(o, t)
	}

	if length == 10 {
		return errors.New("Transaction queue is full")
	}

	return nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (op *Operator) ExecuteTrade(o *Order, tr *Trade) (*types.Transaction, error) {
	tx, err := op.Exchange.Trade(o, t)
	if err != nil {
		return nil, err
	}

	err := op.TradeService.UpdateTradeTx(tr, tx)
	if err != nil {
		return nil, errors.New("Could not update trade tx attribute")
	}

	err := PublishNewTradeExecutedEvent(tr)
	if err != nil {
		return nil, errors.New("Could not publish trade executed event")
	}

	return tx, nil
}


// Validate checks that the operator configuration is sufficient.
func (op *Operator) Validate() error {
	balance, err := op.EthereumClient.getPendingBalanceAt
	if err != nil {
		return err
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


func getQueue(ch *amqp.Channel, queue string) *amqp.Queue {
	if queues[queue] == nil {
		q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %s", err)
		}
		queues[queue] = &q
	}
	return queues[queue]
}

func getChannel(id string) *amqp.Channel {
	if channels[id] == nil {
		ch, err := rabbitmq.Conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
			panic(err)
		}
		channels[id] = ch
	}
	return channels[id]
}

