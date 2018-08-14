package engine

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()
	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	e.addOrder(sampleOrder)
	ssKey, listKey := sampleOrder.GetOBKeys()

	rs, err := s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, sampleOrder.Hash, k, "Expected sorted set value: %v got: %v", sampleOrder.Hash, k)
			assert.Equalf(t, sampleOrder.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", sampleOrder.CreatedAt.Unix(), v)
		}
	}

	rse, err := s.Get(listKey + "::" + sampleOrder.Hash)
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(orderJSON), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, orderJSON, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, sampleOrder.Amount, rse)

}
