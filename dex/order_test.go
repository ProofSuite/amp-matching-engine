package dex

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-test/deep"
)

const testPrivateKey = "7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660"

func TestComputeOrderHash(t *testing.T) {
	order := &Order{
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		AmountBuy:       big.NewInt(1000),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		FeeMake:         big.NewInt(50),
		FeeTake:         big.NewInt(50),
		TokenBuy:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		TokenSell:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Maker:           common.HexToAddress("0xc9b32e9563fe99612ce3a2695ac2a6404c111dde"),
	}

	order.Hash = order.ComputeHash()

	if hash := order.Hash.String(); hash != "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a" {
		t.Error("Expected Orderhash to equal 0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a but got", hash)
	}
}

func TestSignOrder(t *testing.T) {
	o := &Order{
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		AmountBuy:       big.NewInt(1000),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		FeeMake:         big.NewInt(50),
		FeeTake:         big.NewInt(50),
		TokenBuy:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		TokenSell:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Maker:           common.HexToAddress("0xc9b32e9563fe99612ce3a2695ac2a6404c111dde"),
	}

	wallet := NewWalletFromPrivateKey(testPrivateKey)
	o.Sign(wallet)

	expectedSignature := &Signature{
		V: 28,
		R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
		S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
	}

	if o.Signature.V != expectedSignature.V {
		t.Errorf("Expected V to be equal to %v but got %v", expectedSignature.V, o.Signature.V)
	}

	if o.Signature.R != expectedSignature.R {
		t.Errorf("Expected R to be equal to %v but got %v", expectedSignature.R.Hex(), o.Signature.R.Hex())
	}

	if o.Signature.S != expectedSignature.S {
		t.Errorf("Expected S to be equal to %v but got %v", expectedSignature.S.Hex(), o.Signature.S.Hex())
	}
}

func TestVerifySignOrder(t *testing.T) {
	wallet := NewWalletFromPrivateKey(testPrivateKey)

	o := &Order{
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		AmountBuy:       big.NewInt(1000),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		FeeMake:         big.NewInt(50),
		FeeTake:         big.NewInt(50),
		TokenBuy:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		TokenSell:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Maker:           wallet.Address,
	}

	o.Sign(wallet)

	valid, err := o.VerifySignature()
	if err != nil {
		t.Errorf("Error verifying signature: %v", err)
	}

	if !valid {
		t.Errorf("Order signature is not valid")
	}

}

func TestMarshal(t *testing.T) {
	o := &Order{
		Id:              0,
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		Maker:           common.HexToAddress("0xc9b32e9563fe99612ce3a2695ac2a6404c111dde"),
		TokenBuy:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		TokenSell:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		SymbolBuy:       "ZRX",
		SymbolSell:      "WETH",
		AmountBuy:       big.NewInt(1000),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		FeeMake:         big.NewInt(50),
		FeeTake:         big.NewInt(50),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		PairID: common.HexToHash("0x4256b90089bb0cb7c266ec3ac0467845b29797dfdfabb1f2689f71c05f329d5b"),
		Hash:   common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}

	expected := map[string]interface{}{
		"amount":          0,
		"amountBuy":       "1000",
		"amountSell":      "100",
		"exchangeAddress": "0xae55690d4b079460e6ac28aaa58c9ec7b73a7485",
		"expires":         "10000",
		"feeMake":         "50",
		"feeTake":         "50",
		"hash":            "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a",
		"id":              0,
		"maker":           "0xc9b32e9563fe99612ce3a2695ac2a6404c111dde",
		"nonce":           "1000",
		"pairID":          "0x4256b90089bb0cb7c266ec3ac0467845b29797dfdfabb1f2689f71c05f329d5b",
		"price":           0,
		"signature": map[string]interface{}{
			"V": 28,
			"R": "0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85",
			"S": "0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff",
		},
		"symbolBuy":  "ZRX",
		"symbolSell": "WETH",
		"tokenBuy":   "0xe41d2489571d322189246dafa5ebde1f4699f498",
		"tokenSell":  "0x12459c951127e0c374ff9105dda097662a027093",
	}

	encoded, err := json.Marshal(o)
	if err != nil {
		t.Errorf("Error encoding order: %v", err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("Error encoding json: %v", err)
	}

	if diff := deep.Equal(encoded, expectedJSON); diff != nil {
		t.Error(diff)
	}
}

func TestUnMarshal(t *testing.T) {
	expected := Order{
		Id:              0,
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		Maker:           common.HexToAddress("0xc9b32e9563fe99612ce3a2695ac2a6404c111dde"),
		TokenBuy:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		TokenSell:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		SymbolBuy:       "ZRX",
		SymbolSell:      "WETH",
		AmountBuy:       big.NewInt(1000),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		FeeMake:         big.NewInt(50),
		FeeTake:         big.NewInt(50),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		PairID: common.HexToHash("0x4256b90089bb0cb7c266ec3ac0467845b29797dfdfabb1f2689f71c05f329d5b"),
		Hash:   common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}

	payload :=
		`{
			"id": "0",
			"exchangeAddress": "0xae55690d4b079460e6ac28aaa58c9ec7b73a7485",
			"maker":           "0xc9b32e9563fe99612ce3a2695ac2a6404c111dde",
			"tokenBuy":        "0xe41d2489571d322189246dafa5ebde1f4699f498",
			"tokenSell":       "0x12459c951127e0c374ff9105dda097662a027093",
			"symbolBuy":       "ZRX",
			"symbolSell":      "WETH",
			"amountBuy":       "1000",
			"amountSell": "100",
			"expires": "10000",
			"nonce":   "1000",
			"feeMake": "50",
			"feeTake": "50",
			"signature": {
				"V": 28,
				"R": "0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85",
				"S": "0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"
			},
			"pairID": "0x4256b90089bb0cb7c266ec3ac0467845b29797dfdfabb1f2689f71c05f329d5b",
			"hash":   "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"
		}`

	order := Order{}
	err := json.Unmarshal([]byte(payload), &order)
	if err != nil {
		t.Errorf("Could not unmarshal payload: %v", err)
	}

	if diff := deep.Equal(order, expected); diff != nil {
		t.Error(diff)
	}

}
