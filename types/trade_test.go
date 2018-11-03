package types

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestTradeJSON(t *testing.T) {
	expected := &Trade{
		ID:             bson.ObjectIdHex("537f700b537461b70c5f0000"),
		Maker:          common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		Taker:          common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:      common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:     common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Hash:           common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		MakerOrderHash: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		TakerOrderHash: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		PairName:       "ZRX/WETH",
		PricePoint:     big.NewInt(10000),
		Amount:         big.NewInt(100),
	}

	encoded, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("Error encoding order: %v", err)
	}

	trade := &Trade{}
	err = json.Unmarshal([]byte(encoded), &trade)
	if err != nil {
		t.Errorf("Could not unmarshal payload: %v", err)
	}

	if diff := deep.Equal(expected, trade); diff != nil {
		t.Errorf("Expected: \n%+v\nGot: \n%+v\n\n", expected, trade)
	}
}

func TestTradeBSON(t *testing.T) {
	expected := &Trade{
		ID:             bson.ObjectIdHex("537f700b537461b70c5f0000"),
		Maker:          common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		Taker:          common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:      common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:     common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Hash:           common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		MakerOrderHash: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		TakerOrderHash: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		PairName:       "ZRX/WETH",
		PricePoint:     big.NewInt(10000),
		Amount:         big.NewInt(100),
		CreatedAt:      time.Unix(1405544146, 0),
		UpdatedAt:      time.Unix(1405544146, 0),
	}

	data, err := bson.Marshal(expected)
	if err != nil {
		t.Error(err)
	}

	decoded := &Trade{}
	if err := bson.Unmarshal(data, decoded); err != nil {
		t.Error(err)
	}

	assert.Equal(t, decoded, expected)
}
