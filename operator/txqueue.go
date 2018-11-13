package operator

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/streadway/amqp"
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

	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)

	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	err = txq.Broker.ConsumeQueuedTrades(ch, &q, txq.ExecuteTrade)
	if err != nil {
		logger.Error(err)
	}

	return txq, nil
}

func (txq *TxQueue) GetChannel() *amqp.Channel {
	name := "TX_QUEUES" + txq.Name
	return txq.Broker.GetChannel(name)
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

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(m *types.Matches, tag uint64) error {
	logger.Infof("Executing trades: %v", m)

	callOpts := txq.GetTxCallOptions()
	gasLimit, err := txq.Exchange.CallBatchTrades(m, callOpts)
	if err != nil {
		logger.Error(err)
		return err
	}

	if gasLimit < 120000 {
		logger.Warning("GAS LIMIT: ", gasLimit)
		err = txq.Broker.PublishTradeInvalidMessage(m)
		if err != nil {
			logger.Error(err)
			return err
		}

		return errors.New("Invalid Trade")
	}

	nonce, err := txq.EthereumProvider.GetPendingNonceAt(txq.Wallet.Address)
	if err != nil {
		logger.Error(err)
		return err
	}

	txOpts := txq.GetTxSendOptions()
	txOpts.Nonce = big.NewInt(int64(nonce))
	tx, err := txq.Exchange.ExecuteBatchTrades(m, txOpts)
	if err != nil {
		logger.Error(err)
		return err
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
	err = txq.Broker.PublishTradeSentMessage(m)
	if err != nil {
		logger.Error(err)
		return errors.New("Could not update")
	}

	receipt, err := txq.EthereumProvider.WaitMined(tx.Hash())
	if err != nil {
		logger.Error(err)
		return err
	}

	if receipt.Status == 0 {
		logger.Errorf("Reverted transaction: %v", receipt)
		err := txq.HandleTxError(m)
		if err != nil {
			logger.Error(err)
			return err
		}

		return errors.New("Reverted Transaction")
	}

	err = txq.HandleTxSuccess(m, receipt)
	if err != nil {
		logger.Error(err)
		return err
	}

	ch := txq.GetChannel()
	ch.Ack(tag, false)
	return nil
}

func (txq *TxQueue) HandleTxError(m *types.Matches) error {
	logger.Infof("Transaction failed: %v", m)

	errType := "Transaction failed"
	err := txq.Broker.PublishTxErrorMessage(m, errType)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) HandleTxSuccess(m *types.Matches, receipt *eth.Receipt) error {
	logger.Infof("Transaction success: %v", m)

	err := txq.Broker.PublishTradeSuccessMessage(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PublishPendingTrades(m *types.Matches) error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q := txq.Broker.GetQueue(ch, name)

	b, err := json.Marshal(m)
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
