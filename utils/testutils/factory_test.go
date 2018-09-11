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
		BuyToken:        ZRX,
		SellToken:       WETH,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyAmount:       units.Ethers(2),
		SellAmount:      units.Ethers(100),
		FilledAmount:    big.NewInt(0),
		Expires:         big.NewInt(1e18),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Price:           big.NewInt(50), //multiplier from the process order function
		PricePoint:      big.NewInt(50),
		Amount:          units.Ethers(2),
		Side:            "BUY",
		Status:          "NEW",
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
		BuyToken:        WETH,
		SellToken:       ZRX,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyAmount:       units.Ethers(100),
		SellAmount:      units.Ethers(1),
		FilledAmount:    big.NewInt(0),
		Expires:         big.NewInt(1e18),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Side:            "SELL",
		Status:          "NEW",
		PairName:        "ZRX/WETH",
		Nonce:           order.Nonce,
		Signature:       order.Signature,
		Hash:            order.Hash,
		Price:           big.NewInt(100),
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
		BuyToken:        WETH,
		SellToken:       ZRX,
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyAmount:       units.Ethers(2500),
		SellAmount:      units.Ethers(10),
		FilledAmount:    big.NewInt(0),
		Nonce:           order.Nonce,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Expires:         big.NewInt(1e18),
		Signature:       order.Signature,
		Side:            "SELL",
		Status:          "NEW",
		PairName:        "ZRX/WETH",
		Hash:            order.Hash,
		Price:           big.NewInt(250),
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
			Hash: order.Hash.Hex(),
			Data: expectedOrder,
		},
	}

	Compare(t, expectedMessage, msg)
	Compare(t, expectedOrder, order)
}
