package testutils

import (
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

func TestNewOrderFromFactory(t *testing.T) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		t.Errorf("Could not load configuration: %v", err)
	}

	exchangeAddress := common.HexToAddress(app.Config.ExchangeAddress)
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchangeAddress)
	if err != nil {
		t.Errorf("Error creating order factory client: %v", err)
	}

	order, err := f.NewOrder(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	expected := &types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchangeAddress,
		BuyToken:        ZRX,
		SellToken:       WETH,
		BuyAmount:       big.NewInt(1),
		SellAmount:      big.NewInt(1),
		Expires:         order.Expires,
		Nonce:           order.Nonce,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature:       order.Signature,
		Hash:            order.Hash,
		Status:          "NEW",
	}

	Compare(t, expected, order)
}

func TestNewFactoryBuyOrder(t *testing.T) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		t.Errorf("Could not load configuration: %v", err)
	}

	exchangeAddress := common.HexToAddress(app.Config.ExchangeAddress)
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchangeAddress)
	if err != nil {
		t.Errorf("Error creating order factory client: %v", err)
	}

	order, err := f.NewBuyOrder(50, 2)
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	err = order.Process(pair)
	if err != nil {
		t.Errorf("Could not process order: %v", err)
	}

	expected := &types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchangeAddress,
		BuyToken:        ZRX,
		SellToken:       WETH,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyAmount:       big.NewInt(2),
		SellAmount:      big.NewInt(100),
		Expires:         big.NewInt(1e18),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Price:           big.NewInt(50 * 1e8), //multiplier from the process order function
		Amount:          big.NewInt(2),
		Side:            "BUY",
		Status:          "NEW",
		PairName:        "ZRX/WETH",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
	}

	Compare(t, expected, order)
}

func TestNewFactorySellOrder(t *testing.T) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		t.Errorf("Could not load configuration: %v", err)
	}

	exchange := common.HexToAddress(app.Config.ExchangeAddress)
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchange)
	if err != nil {
		t.Errorf("Error creating order factory client: %v", err)
	}

	order, err := f.NewSellOrder(100, 1)
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	err = order.Process(pair)
	if err != nil {
		t.Errorf("Could not process order: %v", err)
	}

	expected := &types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchange,
		BuyToken:        WETH,
		SellToken:       ZRX,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyAmount:       big.NewInt(100),
		SellAmount:      big.NewInt(1),
		Expires:         big.NewInt(1e18),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Side:            "SELL",
		Status:          "NEW",
		PairName:        "ZRX/WETH",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
		Price:           big.NewInt(100 * 1e8),
		Amount:          big.NewInt(1),
	}

	Compare(t, expected, order)
}

func TestNewFactorySellOrder2(t *testing.T) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		t.Errorf("Could not load configuration: %v", err)
	}

	exchange := common.HexToAddress(app.Config.ExchangeAddress)
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchange)
	if err != nil {
		t.Errorf("Error creating factory: %v", err)
	}

	order, err := f.NewSellOrder(250, 10) //Selling 10 ZRX at the price of 1 ZRX = 250 WETH
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	err = order.Process(pair)
	if err != nil {
		t.Errorf("Could not process order: %v", err)
	}

	expected := &types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchange,
		BuyToken:        WETH,
		SellToken:       ZRX,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyAmount:       big.NewInt(2500),
		SellAmount:      big.NewInt(10),
		Nonce:           order.Nonce,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Expires:         big.NewInt(1e18),
		Signature:       order.Signature,
		Side:            "SELL",
		Status:          "NEW",
		PairName:        "ZRX/WETH",
		Hash:            order.Hash,
		Price:           big.NewInt(250 * 1e8),
		Amount:          big.NewInt(10),
	}

	Compare(t, expected, order)
}

func TestNewWebSocketMessage(t *testing.T) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		t.Errorf("Could not load configuration: %v", err)
	}

	exchange := common.HexToAddress(app.Config.ExchangeAddress)
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchange)
	if err != nil {
		t.Errorf("Error creating order factory client: %v", err)
	}

	msg, order, err := f.NewOrderMessage(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Error creating order message: %v", err)
	}

	expectedOrder := &types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchange,
		BuyToken:        ZRX,
		SellToken:       WETH,
		BuyAmount:       big.NewInt(1),
		SellAmount:      big.NewInt(1),
		Expires:         big.NewInt(1e18),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Status:          "NEW",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
	}

	expectedMessage := &types.WebSocketMessage{
		Channel: "orders",
		Payload: types.WebSocketPayload{
			Type: "NEW_ORDER",
			Hash: "",
			Data: expectedOrder,
		},
	}

	Compare(t, expectedMessage, msg)
	Compare(t, expectedOrder, order)
}
