package operator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

var logger = utils.OperatorLogger

// Operator manages the transaction queue that will eventually be
// sent to the exchange contract. The Operator Wallet must be equal to the
// account that initially deployed the exchange contract or an address with operator rights
// on the contract
type Operator struct {
	// AccountService     interfaces.AccountService
	WalletService      interfaces.WalletService
	TradeService       interfaces.TradeService
	OrderService       interfaces.OrderService
	EthereumProvider   interfaces.EthereumProvider
	Exchange           interfaces.Exchange
	TxQueues           []*TxQueue
	QueueAddressIndex  map[common.Address]*TxQueue
	RabbitMQConnection *rabbitmq.Connection
	mutex              *sync.Mutex
}

type OperatorInterface interface {
	SubscribeOperatorMessages(fn func(*types.OperatorMessage) error) error
	QueueTrade(o *types.Order, t *types.Trade) error
	GetShortestQueue() (*TxQueue, int, error)
	SetFeeAccount(account common.Address) (*eth.Transaction, error)
	SetOperator(account common.Address, isOperator bool) (*eth.Transaction, error)
	FeeAccount() (common.Address, error)
	Operator(addr common.Address) (bool, error)
}

// NewOperator creates a new operator struct. It creates an exchange contract instance from the
// provided address. The error and trade events are received in the ErrorChannel and TradeChannel.
// Upon receiving errors and trades in their respective channels, event payloads are sent to the
// associated order maker and taker sockets through the through the event channel on the Order and Trade struct.
// In addition, an error event cancels the trade in the trading engine and makes the order available again.
func NewOperator(
	walletService interfaces.WalletService,
	tradeService interfaces.TradeService,
	orderService interfaces.OrderService,
	provider interfaces.EthereumProvider,
	exchange interfaces.Exchange,
	conn *rabbitmq.Connection,
) (*Operator, error) {
	txqueues := []*TxQueue{}
	addressIndex := make(map[common.Address]*TxQueue)

	wallets, err := walletService.GetOperatorWallets()
	if err != nil {
		panic(err)
	}

	for i, w := range wallets {
		name := strconv.Itoa(i) + w.Address.Hex()
		ch := conn.GetChannel("TX_QUEUES:" + name)
		err := conn.DeclareQueue(ch, "TX_QUEUES:"+name)
		if err != nil {
			panic(err)
		}

		txq := &TxQueue{
			Name:             name,
			Wallet:           w,
			TradeService:     tradeService,
			EthereumProvider: provider,
			Exchange:         exchange,
			RabbitMQConn:     conn,
		}

		txqueues = append(txqueues, txq)
	}

	op := &Operator{
		WalletService:     walletService,
		TradeService:      tradeService,
		OrderService:      orderService,
		EthereumProvider:  provider,
		Exchange:          exchange,
		TxQueues:          txqueues,
		QueueAddressIndex: addressIndex,
		mutex:             &sync.Mutex{},
	}

	// err = op.PurgeQueues()
	// if err != nil {
	// 	panic(err)
	// }

	go op.HandleEvents()
	return op, nil
}

// SubscribeOperatorMessages
func (op *Operator) SubscribeOperatorMessages(fn func(*types.OperatorMessage) error) error {
	ch := op.RabbitMQConnection.GetChannel("OPERATOR_SUB")
	q := op.RabbitMQConnection.GetQueue(ch, "TX_MESSAGES")

	go func() {
		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			logger.Error("Failed to register a consumer", err)
		}

		forever := make(chan bool)

		go func() {
			for m := range msgs {
				om := &types.OperatorMessage{}
				err := json.Unmarshal(m.Body, &om)
				if err != nil {
					logger.Infof("Error: %v", err)
					continue
				}
				go fn(om)
			}
		}()

		<-forever
	}()
	return nil
}

// Bug: In certain cases, the trade channel seems to be receiving additional unexpected trades.
// In the case TestSocketExecuteOrder (in file socket_test.go) is run on its own, everything is working correctly.
// However, in the case TestSocketExecuteOrder is run among other tests, some tradeLogs do not correspond to an
// order hash in the ordertrade mapping. I suspect this is because the event listener catches events from previous
// tests. It might be helpful to see how to listen to events from up to a certain block.
func (op *Operator) HandleEvents() error {
	tradeEvents, err := op.Exchange.ListenToTrades()
	if err != nil {
		logger.Error(err)
		return err
	}

	errorEvents, err := op.Exchange.ListenToErrors()
	if err != nil {
		logger.Error(err)
		return err
	}

	for {
		select {
		case event := <-errorEvents:
			fmt.Println("TRADE_ERROR_EVENT")
			tradeHash := event.TradeHash
			errID := int(event.ErrorId)

			logger.Info("The error ID is: ", errID)

			tr, err := op.TradeService.GetByHash(tradeHash)
			if err != nil {
				logger.Error(err)
			}

			or, err := op.OrderService.GetByHash(tr.OrderHash)
			if err != nil {
				logger.Error(err)
			}

			go func() {
				err = op.RabbitMQConnection.PublishTxErrorMessage(tr, errID)
				if err != nil {
					logger.Error(err)
				}

				err = op.RabbitMQConnection.PublishTradeCancelMessage(or, tr)
				if err != nil {
					logger.Error(err)
				}
			}()

		case event := <-tradeEvents:
			tr, err := op.TradeService.GetByHash(event.TradeHash)
			if err != nil {
				logger.Error(err)
			}

			logger.Info("TRADE_SUCCESS_EVENT", tr.Hash.Hex(), tr.TxHash.Hex())

			or, err := op.OrderService.GetByHash(tr.OrderHash)
			if err != nil {
				logger.Error(err)
			}

			go func() {
				_, err := op.EthereumProvider.WaitMined(tr.TxHash)
				if err != nil {
					logger.Error(err)
				}

				logger.Info("TRADE_MINED IN HANDLE EVENTS: ", tr.Hash.Hex())

				err = op.RabbitMQConnection.PublishTradeSuccessMessage(or, tr)
				if err != nil {
					logger.Error(err)
				}
			}()
		}
	}
}

func (op *Operator) HandleTrades(msg *types.OperatorMessage) error {
	o := msg.Order
	// t := msg.Trade

	//TODO move this to the order service
	err := o.Validate()
	if err != nil {
		logger.Error(err)
		return err
	}

	//TODO move this to the order service
	ok, err := o.VerifySignature()
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Invalid signature")
	}

	err = op.QueueTrade(msg.Order, msg.Trade)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// QueueTrade
func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	txq, len, err := op.GetShortestQueue()
	if err != nil {
		logger.Error(err)
		return err
	}

	if len > 10 {
		logger.Info("Transaction queue is full")
		return errors.New("Transaction queue is full")
	}

	logger.Info("QUEING TRADE", len)
	err = txq.QueueTrade(o, t)
	if err != nil {
		logger.Warning("INVALID TRADE")
		return err
	}

	return nil
}

// GetShortestQueue
func (op *Operator) GetShortestQueue() (*TxQueue, int, error) {
	shortest := &TxQueue{}
	min := 1000

	for _, txq := range op.TxQueues {
		if shortest == nil {
			shortest = txq
			min = txq.Length()
		}

		ln := txq.Length()
		if ln < min {
			shortest = txq
			min = ln
		}
	}

	return shortest, min, nil
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (op *Operator) SetFeeAccount(account common.Address) (*eth.Transaction, error) {
	txOpts, err := op.GetTxSendOptions()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tx, err := op.Exchange.SetFeeAccount(account, txOpts)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (op *Operator) SetOperator(account common.Address, isOperator bool) (*eth.Transaction, error) {
	txOpts, err := op.GetTxSendOptions()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tx, err := op.Exchange.SetOperator(account, isOperator, txOpts)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tx, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (op *Operator) FeeAccount() (common.Address, error) {
	account, err := op.Exchange.FeeAccount()
	if err != nil {
		logger.Error(err)
		return common.Address{}, err
	}

	return account, nil
}

// Operator returns true if the given address is an operator of the exchange and returns false otherwise
func (op *Operator) Operator(addr common.Address) (bool, error) {
	isOperator, err := op.Exchange.Operator(addr)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	return isOperator, nil
}

func (op *Operator) PurgeQueues() error {
	for _, txq := range op.TxQueues {
		err := txq.PurgePendingTrades()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (op *Operator) GetTxSendOptions() (*bind.TransactOpts, error) {
	wallet, err := op.WalletService.GetDefaultAdminWallet()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return bind.NewKeyedTransactor(wallet.PrivateKey), nil
}

// func (op *Operator) ValidateTrade(o *types.Order, t *types.Trade) error {
// 	// fee balance validation
// 	wethAddress := common.HexToAddress(app.Config.Ethereum["weth_address"])
// 	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])

// 	makerBalanceRecord, err := op.AccountService.GetTokenBalances(o.UserAddress)
// 	if err != nil {
// 		logger.Error("Error retrieving maker token balances", err)
// 		return err
// 	}

// 	takerBalanceRecord, err := op.AccountService.GetTokenBalances(t.Taker)
// 	if err != nil {
// 		logger.Error("Error retrieving taker token balances", err)
// 		return err
// 	}

// 	makerWethBalance, err := op.EthereumProvider.BalanceOf(o.UserAddress, wethAddress)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	makerWethAllowance, err := op.EthereumProvider.Allowance(o.UserAddress, exchangeAddress, wethAddress)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	makerTokenBalance, err := op.EthereumProvider.BalanceOf(o.UserAddress, o.SellToken)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	makerTokenAllowance, err := op.EthereumProvider.Allowance(o.UserAddress, exchangeAddress, o.SellToken)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	takerWethBalance, err := op.EthereumProvider.BalanceOf(t.Taker, wethAddress)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	takerWethAllowance, err := op.EthereumProvider.Allowance(t.Taker, exchangeAddress, wethAddress)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	takerTokenBalance, err := op.EthereumProvider.BalanceOf(t.Taker, o.BuyToken)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	takerTokenAllowance, err := op.EthereumProvider.Allowance(t.Taker, exchangeAddress, o.BuyToken)
// 	if err != nil {
// 		logger.Error("Error", err)
// 		return err
// 	}

// 	fee := math.Max(o.MakeFee, o.TakeFee)
// 	makerAvailableWethBalance := math.Sub(makerWethBalance, makerBalanceRecord[wethAddress].LockedBalance)
// 	makerAvailableTokenBalance := math.Sub(makerTokenBalance, makerBalanceRecord[o.SellToken].LockedBalance)
// 	takerAvailableWethBalance := math.Sub(takerWethBalance, takerBalanceRecord[wethAddress].LockedBalance)
// 	takerAvailableTokenBalance := math.Sub(takerTokenBalance, takerBalanceRecord[o.BuyToken].LockedBalance)

// 	if makerAvailableWethBalance.Cmp(fee) == -1 {
// 		return errors.New("Insufficient WETH Balance")
// 	}

// 	if makerWethAllowance.Cmp(fee) == -1 {
// 		return errors.New("Insufficient WETH Balance")
// 	}

// 	if makerAvailableSellTokenBalance.Cmp(o.SellAmount) != 1 {
// 		return errors.New("Insufficient Balance")
// 	}

// 	if makerTokenAllowance.Cmp(o.SellAmount) != 1 {
// 		return errors.New("Insufficient Allowance")
// 	}

// 	if takerAvailableWethBalance.Cmp(fee) == -1 {
// 		return errors.New("Insufficient WETH Balance")
// 	}

// 	if takerWethAllowance.Cmp(fee) == -1 {
// 		return errors.New("Insufficient WETH Balance")
// 	}

// 	if takerAvailableTokenBalance.Cmp(t.Amount) != 1 {
// 		return errors.New("Insufficient Balance")
// 	}

// 	if takerTokenAllowance.Cmp(t.Amount) != 1 {
// 		return errors.New("Insufficient Allowance")
// 	}

// 	return nil
// }
