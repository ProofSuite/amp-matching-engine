package testutils

import (
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/units"
)

func TestNewOrderFromFactory(t *testing.T) {
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	exchangeAddress := GetTestAddress2()
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
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		Amount:          big.NewInt(1),
		Nonce:           order.Nonce,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature:       order.Signature,
		Hash:            order.Hash,
		Status:          "OPEN",
	}

	Compare(t, expected, order)
}

func TestNewFactoryBuyOrder(t *testing.T) {
	exchangeAddress := GetTestAddress3()
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

	expected := types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchangeAddress,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		FilledAmount:    big.NewInt(0),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		PricePoint:      big.NewInt(50),
		Amount:          units.Ethers(2),
		Side:            "BUY",
		Status:          "OPEN",
		PairName:        "ZRX/WETH",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
	}

	CompareOrder(t, &expected, &order)
}

func TestNewFactorySellOrder1(t *testing.T) {
	exchangeAddress := GetTestAddress3()
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchangeAddress)
	if err != nil {
		t.Errorf("Error creating order factory client: %v", err)
	}

	order, err := f.NewSellOrder(100, 1)
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	expected := types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchangeAddress,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		FilledAmount:    big.NewInt(0),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Side:            "SELL",
		Status:          "OPEN",
		PairName:        "ZRX/WETH",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
		PricePoint:      big.NewInt(100),
		Amount:          units.Ethers(1),
	}

	CompareOrder(t, &expected, &order)
}

func TestNewFactorySellOrder2(t *testing.T) {
	exchangeAddress := GetTestAddress3()
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchangeAddress)
	if err != nil {
		t.Errorf("Error creating factory: %v", err)
	}

	order, err := f.NewSellOrder(250, 10) //Selling 10 ZRX at the price of 1 ZRX = 250 WETH
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	expected := types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchangeAddress,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		FilledAmount:    big.NewInt(0),
		Nonce:           order.Nonce,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature:       order.Signature,
		Side:            "SELL",
		Status:          "OPEN",
		PairName:        "ZRX/WETH",
		Hash:            order.Hash,
		PricePoint:      big.NewInt(250),
		Amount:          units.Ethers(10),
	}

	CompareOrder(t, &expected, &order)
}

func TestNewWebSocketMessage(t *testing.T) {
	exchangeAddress := GetTestAddress3()
	pair := GetZRXWETHTestPair()
	wallet := GetTestWallet1()
	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	f, err := NewOrderFactory(pair, wallet, exchangeAddress)
	if err != nil {
		t.Errorf("Error creating order factory client: %v", err)
	}

	msg, order, err := f.NewOrderMessage(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Error creating order message: %v", err)
	}

	expectedOrder := &types.Order{
		UserAddress:     wallet.Address,
		ExchangeAddress: exchangeAddress,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		Amount:          big.NewInt(1),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Status:          "OPEN",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
	}

	expectedMessage := &types.WebsocketMessage{
		Channel: "orders",
		Event: types.WebsocketEvent{
			Type:    "NEW_ORDER",
			Hash:    order.Hash.Hex(),
			Payload: expectedOrder,
		},
	}

	Compare(t, expectedMessage, msg)
	Compare(t, expectedOrder, order)
}
