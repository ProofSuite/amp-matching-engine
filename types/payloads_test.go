package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-test/deep"
)

func TestNewOrderPayload(t *testing.T) {
	p := &NewOrderPayload{
		PairName:        "ZRX/WETH",
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Amount:          big.NewInt(1000),
		PricePoint:      big.NewInt(100),
		MakeFee:         big.NewInt(50),
		TakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash: common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}

	encoded, err := json.Marshal(p)
	if err != nil {
		t.Errorf("Error encoding order: %v", err)
	}

	decoded := &NewOrderPayload{}
	err = json.Unmarshal([]byte(encoded), decoded)
	if err != nil {
		t.Errorf("Could not unmarshal payload: %v", err)
	}

	if diff := deep.Equal(p, decoded); diff != nil {
		fmt.Printf("Expected: \n%+v\nGot: \n%+v\n\n", p, decoded)
	}
}
