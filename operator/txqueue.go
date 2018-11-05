package operator

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth "github.com/ethereum/go-ethereum/core/types"
)

type TxQueue struct {
	Name             string
	Wallet           *types.Wallet
	TradeService     interfaces.TradeService
	OrderService     interfaces.OrderService
	EthereumProvider interfaces.EthereumProvider
	Exchange         interfaces.Exchange
	Broker           *rabbitmq.Connection
}

// NewTxQueue
func NewTxQueue(
	n string,
	tr interfaces.TradeService,
	p interfaces.EthereumProvider,
	o interfaces.OrderService,
	w *types.Wallet,
	ex interfaces.Exchange,
	rabbitConn *rabbitmq.Connection,
) (*TxQueue, error) {
	txq := &TxQueue{
		Name:             n,
		TradeService:     tr,
		OrderService:     o,
		EthereumProvider: p,
		Wallet:           w,
		Exchange:         ex,
		Broker:           rabbitConn,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return txq, nil
}

func (txq *TxQueue) GetTxSendOptions() *bind.TransactOpts {
	return bind.NewKeyedTransactor(txq.Wallet.PrivateKey)
}

func (txq *TxQueue) GetTxCallOptions() *ethereum.CallMsg {
	address := txq.Exchange.GetAddress()

	return &ethereum.CallMsg{From: txq.Wallet.Address, To: &address}
}

// Length
func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	return q.Messages
}

// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
// func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
// TODO: Currently doesn't seem thread safe and fails unless called with a sleep time between each call.
func (txq *TxQueue) QueueTrade(m *types.Matches) error {
	logger.Info("QUEUE LENGTH", txq.Length())
	if txq.Length() == 0 {
		_, err := txq.ExecuteTrade(m)
		if err != nil {
			logger.Error(err)
			logger.Info("This is an invalid trade")
			return err
		}
	}

	err := txq.PublishPendingTrades(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(m *types.Matches) (*eth.Transaction, error) {
	logger.Info("EXECUTING TRADE", m)

	callOpts := txq.GetTxCallOptions()
	gasLimit, err := txq.Exchange.CallBatchTrades(m, callOpts)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if gasLimit < 120000 {
		logger.Warning("GAS LIMIT: ", gasLimit)
		err = txq.Broker.PublishTradeInvalidMessage(m)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		go txq.ExecuteNextTrade()
		return nil, errors.New("Invalid Trade")
	}

	nonce, err := txq.EthereumProvider.GetPendingNonceAt(txq.Wallet.Address)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	txOpts := txq.GetTxSendOptions()
	txOpts.Nonce = big.NewInt(int64(nonce))
	tx, err := txq.Exchange.ExecuteBatchTrades(m, txOpts)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	updatedTrades := []*types.Trade{}
	for _, t := range m.Trades {
		updated, err := txq.TradeService.UpdatePendingTrade(t, tx.Hash())
		if err != nil {
			logger.Error(err)
		}

		updatedTrades = append(updatedTrades, updated)
	}

	m.Trades = updatedTrades

	utils.PrintJSON(m.Trades)

	err = txq.Broker.PublishTradeSentMessage(m)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Could not update")
	}

	go func() {
		_, err := txq.EthereumProvider.WaitMined(tx.Hash())
		if err != nil {
			logger.Error(err)
		}

		// logger.Info("TRADE_MINED IN EXECUTE TRADE: ", tr.Hash.Hex())
		//TODO in this case, what happens in the case we have a lot of trades, i think the best solution
		//TODO is to register the events in redis, and check if they already exist.
		len := txq.Length()
		if len > 0 {
			nextMatch, err := txq.PopPendingTrades()
			if err != nil {
				logger.Error(err)
				return
			}

			if nextMatch == nil {
				return
			}

			// very hacky
			if nextMatch.TakerOrderHash() == m.TakerOrderHash() {
				nextMatch, err = txq.PopPendingTrades()
				if err != nil {
					logger.Error(err)
					return
				}

				if nextMatch == nil {
					return
				}
			}

			go txq.ExecuteTrade(nextMatch)
		}
	}()

	return tx, nil
}

func (txq *TxQueue) ExecuteNextTrade() error {
	len := txq.Length()
	logger.Info("LENGTH of the queue is ", len)
	if len > 0 {
		match, err := txq.PopPendingTrades()
		if err != nil {
			logger.Error(err)
			return err
		}

		// logger.Info("NEXT_TRADE: ", msg.Trade.Hash.Hex())
		go txq.ExecuteTrade(match)
		return nil
	}

	return nil
}

func (txq *TxQueue) PublishPendingTrades(m *types.Matches) error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q := txq.Broker.GetQueue(ch, name)

	msg := &types.PendingTradeBatch{m}
	b, err := json.Marshal(msg)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = txq.Broker.Publish(ch, q, b)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PurgePendingTrades() error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)

	err := txq.Broker.Purge(ch, name)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PopPendingTrade
func (txq *TxQueue) PopPendingTrades() (*types.Matches, error) {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q := txq.Broker.GetQueue(ch, name)

	msg, _, _ := ch.Get(
		q.Name,
		true,
	)

	if len(msg.Body) == 0 {
		return nil, nil
	}

	pding := &types.Matches{}
	err := json.Unmarshal(msg.Body, &pding)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = pding.Validate()
	if err != nil {
		return nil, nil
	}

	return pding, nil
}
