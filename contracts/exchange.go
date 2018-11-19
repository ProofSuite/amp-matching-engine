package contracts

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

type ethereumClientInterface interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	// PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *eth.Transaction) error
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error)
}

// Exchange is an augmented interface to the Exchange.sol smart-contract. It uses the
// smart-contract bindings generated with abigen and adds additional functionality and
// simplifications to these bindings.
// Address is the Ethereum address of the exchange contract
// Contract is the original abigen bindings
// CallOptions are options for making read calls to the connected backend
// TxOptions are options for making write txs to the connected backend
type Exchange struct {
	Address       common.Address
	WalletService interfaces.WalletService
	Interface     *contractsinterfaces.Exchange
	Client        ethereumClientInterface
}

// Returns a new exchange interface for a given wallet, contract address and connected backend.
// The exchange contract need to be already deployed at the given address. The given wallet will
// be used by default when sending transactions with this object.
func NewExchange(
	w interfaces.WalletService,
	contractAddress common.Address,
	backend ethereumClientInterface,
) (*Exchange, error) {
	instance, err := contractsinterfaces.NewExchange(contractAddress, backend)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &Exchange{
		WalletService: w,
		Interface:     instance,
		Client:        backend,
		Address:       contractAddress,
	}, nil
}

func (e *Exchange) GetAddress() common.Address {
	return e.Address
}

func (e *Exchange) DefaultTxOptions() (*bind.TransactOpts, error) {
	wallet, err := e.WalletService.GetDefaultAdminWallet()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	opts := bind.NewKeyedTransactor(wallet.PrivateKey)
	return opts, nil
}

func (e *Exchange) GetTxCallOptions() *bind.CallOpts {
	return &bind.CallOpts{Pending: true}
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (e *Exchange) SetFeeAccount(a common.Address, txOpts *bind.TransactOpts) (*eth.Transaction, error) {
	tx, err := e.Interface.SetFeeAccount(txOpts, a)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (e *Exchange) SetOperator(a common.Address, isOperator bool, txOpts *bind.TransactOpts) (*eth.Transaction, error) {
	tx, err := e.Interface.SetOperator(txOpts, a, isOperator)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tx, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (e *Exchange) FeeAccount() (common.Address, error) {
	callOptions := e.GetTxCallOptions()

	account, err := e.Interface.RewardAccount(callOptions)
	if err != nil {
		logger.Error(err)
		return common.Address{}, err
	}

	return account, nil
}

func (e *Exchange) Operator(a common.Address) (bool, error) {
	callOptions := e.GetTxCallOptions()
	// Operator returns true if the given address is an operator of the exchange and returns false otherwise
	isOperator, err := e.Interface.Operators(callOptions, a)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	return isOperator, nil
}

func (e *Exchange) ExecuteBatchTrades(matches *types.Matches, txOpts *bind.TransactOpts) (*eth.Transaction, error) {
	orderValues := [][10]*big.Int{}
	orderAddresses := [][4]common.Address{}
	vValues := [][2]uint8{}
	rsValues := [][4][32]byte{}
	amounts := []*big.Int{}

	makerOrders := matches.MakerOrders
	trades := matches.Trades
	takerOrder := matches.TakerOrder

	for i, _ := range makerOrders {
		mo := makerOrders[i]
		to := takerOrder
		t := trades[i]

		orderValues = append(orderValues, [10]*big.Int{mo.Amount, mo.PricePoint, mo.EncodedSide(), mo.Nonce, to.Amount, to.PricePoint, to.EncodedSide(), to.Nonce, mo.MakeFee, mo.TakeFee})
		orderAddresses = append(orderAddresses, [4]common.Address{mo.UserAddress, to.UserAddress, mo.BaseToken, to.QuoteToken})
		vValues = append(vValues, [2]uint8{mo.Signature.V, to.Signature.V})
		rsValues = append(rsValues, [4][32]byte{mo.Signature.R, mo.Signature.S, to.Signature.R, to.Signature.S})
		amounts = append(amounts, t.Amount)
	}

	tx, err := e.Interface.ExecuteBatchTrades(txOpts, orderValues, orderAddresses, amounts, vValues, rsValues)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tx, nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (e *Exchange) Trade(match *types.Matches, txOpts *bind.TransactOpts) (*eth.Transaction, error) {
	mo := match.MakerOrders[0]
	to := match.TakerOrder
	t := match.Trades[0]

	orderValues := [10]*big.Int{mo.Amount, mo.PricePoint, mo.EncodedSide(), mo.Nonce, to.Amount, to.PricePoint, to.EncodedSide(), to.Nonce, mo.MakeFee, mo.TakeFee}
	orderAddresses := [4]common.Address{mo.UserAddress, to.UserAddress, mo.BaseToken, to.QuoteToken}
	vValues := [2]uint8{mo.Signature.V, to.Signature.V}
	rsValues := [4][32]byte{mo.Signature.R, mo.Signature.S, to.Signature.R, to.Signature.S}
	amount := t.Amount

	tx, err := e.Interface.ExecuteSingleTrade(txOpts, orderValues, orderAddresses, amount, vValues, rsValues)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tx, nil
}

func (e *Exchange) CallBatchTrades(matches *types.Matches, call *ethereum.CallMsg) (uint64, error) {
	orderValues := [][10]*big.Int{}
	orderAddresses := [][4]common.Address{}
	amounts := []*big.Int{}
	vValues := [][2]uint8{}
	rsValues := [][4][32]byte{}

	makerOrders := matches.MakerOrders
	trades := matches.Trades
	takerOrder := matches.TakerOrder

	for i, _ := range makerOrders {
		mo := makerOrders[i]
		to := takerOrder
		t := trades[i]

		orderValues = append(orderValues, [10]*big.Int{mo.Amount, mo.PricePoint, mo.EncodedSide(), mo.Nonce, to.Amount, to.PricePoint, to.EncodedSide(), to.Nonce, mo.MakeFee, mo.TakeFee})
		orderAddresses = append(orderAddresses, [4]common.Address{mo.UserAddress, to.UserAddress, mo.BaseToken, mo.QuoteToken})
		amounts = append(amounts, t.Amount)

		if mo.Signature == nil {
			return 0, errors.New("Maker order is not signed")
		}

		if to.Signature == nil {
			return 0, errors.New("Taker order is not signed")
		}

		vValues = append(vValues, [2]uint8{mo.Signature.V, to.Signature.V})
		rsValues = append(rsValues, [4][32]byte{mo.Signature.R, mo.Signature.S, to.Signature.R, to.Signature.S})
	}

	exchangeABI, err := abi.JSON(strings.NewReader(contractsinterfaces.ExchangeABI))
	if err != nil {
		return 0, err
	}

	data, err := exchangeABI.Pack("executeBatchTrades", orderValues, orderAddresses, amounts, vValues, rsValues)
	if err != nil {
		return 0, err
	}

	// call.Data = data
	// b, err := e.Client.PendingCallContract(context.Background(), *call)
	// if err != nil {
	// 	logger.Error(err)
	// 	return 0, err
	// }

	call.Data = data
	gasLimit, err := e.Client.(bind.ContractBackend).EstimateGas(context.Background(), *call)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return gasLimit, nil
}

func (e *Exchange) CallTrade(match *types.Matches, call *ethereum.CallMsg) (uint64, error) {
	mo := match.MakerOrders[0]
	to := match.TakerOrder
	t := match.Trades[0]

	orderValues := [10]*big.Int{mo.Amount, mo.PricePoint, mo.EncodedSide(), mo.Nonce, to.Amount, to.PricePoint, to.EncodedSide(), to.Nonce, mo.MakeFee, mo.TakeFee}
	orderAddresses := [4]common.Address{mo.UserAddress, to.UserAddress, mo.BaseToken, to.QuoteToken}
	vValues := [2]uint8{mo.Signature.V, to.Signature.V}
	rsValues := [4][32]byte{mo.Signature.R, mo.Signature.S, to.Signature.R, to.Signature.S}

	exchangeABI, err := abi.JSON(strings.NewReader(contractsinterfaces.ExchangeABI))
	if err != nil {
		return 0, err
	}

	data, err := exchangeABI.Pack("executeSingleTrade", orderValues, orderAddresses, t.Amount, vValues, rsValues)
	if err != nil {
		return 0, err
	}

	call.Data = data
	gasLimit, err := e.Client.(bind.ContractBackend).EstimateGas(context.Background(), *call)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return gasLimit, nil
}

// ListenToErrorEvents returns a channel that receives errors logs (events) from the exchange smart contract.
// The error IDs correspond to the following codes:
// 1. MAKER_INSUFFICIENT_BALANCE,
// 2. TAKER_INSUFFICIENT_BALANCE,
// 3. WITHDRAW_INSUFFICIENT_BALANCE,
// 4. WITHDRAW_FEE_TO_HIGH,
// 5. ORDER_EXPIRED,
// 6. WITHDRAW_ALREADY_COMPLETED,
// 7. TRADE_ALREADY_COMPLETED,
// 8. TRADE_AMOUNT_TOO_BIG,
// 9. SIGNATURE_INVALID,
// 10. MAKER_SIGNATURE_INVALID,
// 11. TAKER_SIGNATURE_INVALID
func (e *Exchange) ListenToErrors() (chan *contractsinterfaces.ExchangeLogError, error) {
	events := make(chan *contractsinterfaces.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, events)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return events, nil
}

// ListenToTrades returns a channel that receivs trade logs (events) from the underlying exchange smart contract
func (e *Exchange) ListenToTrades() (chan *contractsinterfaces.ExchangeLogTrade, error) {
	events := make(chan *contractsinterfaces.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, events, nil, nil, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return events, nil
}

func (e *Exchange) ListenToBatchTrades() (chan *contractsinterfaces.ExchangeLogBatchTrades, error) {
	events := make(chan *contractsinterfaces.ExchangeLogBatchTrades)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogBatchTrades(opts, events, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return events, nil
}

func (e *Exchange) GetErrorEvents(logs chan *contractsinterfaces.ExchangeLogError) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, logs)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Exchange) GetTrades(logs chan *contractsinterfaces.ExchangeLogTrade) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, logs, nil, nil, nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Exchange) PrintTrades() error {
	events := make(chan *contractsinterfaces.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, events, nil, nil, nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	go func() {
		for {
			event := <-events
			logger.Infof("New event: %v", event)
		}
	}()

	return nil
}

func (e *Exchange) PrintErrors() error {
	events := make(chan *contractsinterfaces.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, events)
	if err != nil {
		return err
	}

	go func() {
		for {
			event := <-events
			logger.Warningf("New Error Event. Id: %v, Hash: %v\n\n", event.ErrorId, hex.EncodeToString(event.TakerOrderHash[:]))
		}
	}()

	return nil
}

func PrintErrorLog(log *contractsinterfaces.ExchangeLogError) string {
	return fmt.Sprintf("Error:\nErrorID: %v\nOrderHash: %v\n\n", log.ErrorId, log.TakerOrderHash)
}

func PrintTradeLog(log *contractsinterfaces.ExchangeLogTrade) string {
	return fmt.Sprintf("Error:\nMaker: %v\nTaker: %v\nTokenBuy: %v\nTokenSell: %v\nOrderHash: %v\nTradeHash: %v\n\n",
		log.Maker, log.Taker, log.TokenBuy, log.TokenSell, log.OrderHash, log.TradeHash)
}

func PrintCancelOrderLog(log *contractsinterfaces.ExchangeLogCancelOrder) string {
	return fmt.Sprintf("Error:\nSender: %v\nOrderHash: %v\n\n", log.UserAddress, log.OrderHash)
}
