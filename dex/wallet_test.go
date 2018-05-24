package dex

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()

	if reflect.TypeOf(*wallet) != reflect.TypeOf(Wallet{}) {
		t.Error("Wallet type is not correct")
	}

	address := wallet.GetAddress()
	if addressLength := len(address); addressLength != 42 {
		t.Error("Expected address length to be 40, but got: ", addressLength)
	}

	privateKey := wallet.GetPrivateKey()
	if privateKeyLength := len(privateKey); privateKeyLength != 64 {
		t.Error("Expected private key length to be 64, but got: ", privateKeyLength)
	}
}

func TestNewWalletFromPrivateKey(t *testing.T) {
	key := "7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660"

	wallet := NewWalletFromPrivateKey(key)
	if address := wallet.GetAddress(); address != "0xE8E84ee367BC63ddB38d3D01bCCEF106c194dc47" {
		t.Error("Expected address to equal 0xE8E84ee367BC63ddB38d3D01bCCEF106c194dc47 but got: ", address)
	}
}

func TestNewOrder(t *testing.T) {
	wallet := config.Wallets[0]
	pair := config.TokenPairs["ZRXWETH"]

	order, err := wallet.NewOrder(0, 100, 100, pair, BUY)
	if err != nil {
		t.Errorf("Error in function NewOrder: %v", err)
	}

	if order.Id != 0 {
		t.Errorf("Expected ID to be equal to %v but got %v instead", 0, order.Id)
	}

	if order.AmountBuy.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected ID to be equal to %v but got %v instead", 100, order.AmountBuy)
	}

	if order.AmountSell.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected ID to be equal to %v but got %v instead", 100, order.AmountSell)
	}

	if order.TokenBuy != pair.QuoteToken.Address {
		t.Errorf("Expected Token Buy Address to be equal to %v but got %v instead", pair.QuoteToken.Address, order.TokenBuy)
	}

	if order.TokenSell != pair.BaseToken.Address {
		t.Errorf("Expected Token Sell Address to be equal to %v but got %v instead", pair.BaseToken.Address, order.TokenSell)
	}

	if order.SymbolBuy != pair.QuoteToken.Symbol {
		t.Errorf("Expected Token Buy Symbol to be equal to %v but got %v instead", pair.QuoteToken.Symbol, order.SymbolBuy)
	}

	if order.SymbolSell != pair.BaseToken.Symbol {
		t.Errorf("Expected Token Sell Symbol to be equal to %v but got %v instead", pair.BaseToken.Symbol, order.SymbolSell)
	}

	if order.Maker != wallet.Address {
		t.Errorf("Expected Order Maker to be equal to %v but got %v instead", wallet.Address, order.Maker)
	}
}
