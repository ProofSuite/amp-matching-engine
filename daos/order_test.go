package daos

import (
	"math/big"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/units"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	// temp, _ := ioutil.TempDir("", "test")
	// server.SetPath(temp)

	// session := server.Session()
	// db = &Database{session}
}

func TestUpdateOrderByHash(t *testing.T) {
	o := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		PricePoint:      big.NewInt(1000),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "OPEN",
		Side:            "BUY",
		PairName:        "ZRX/WETH",
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	dao := NewOrderDao()

	err := dao.Create(o)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	updated := &types.Order{
		ID:           o.ID,
		UserAddress:  o.UserAddress,
		BaseToken:    o.BaseToken,
		QuoteToken:   o.QuoteToken,
		PricePoint:   big.NewInt(4000),
		Amount:       big.NewInt(4000),
		FilledAmount: big.NewInt(200),
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Signature:    o.Signature,
		Hash:         o.Hash,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}

	err = dao.UpdateByHash(
		o.Hash,
		updated,
	)

	if err != nil {
		t.Errorf("Could not updated order from hash %v", err)
	}

	queried, err := dao.GetByHash(o.Hash)
	if err != nil {
		t.Errorf("Could not get order by hash")
	}

	testutils.CompareOrder(t, updated, queried)
}

func TestOrderUpdate(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Errorf("Could not drop previous order state")
	}

	o := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:  common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		BaseToken:    common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:   common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		PricePoint:   big.NewInt(1000),
		Amount:       big.NewInt(1000),
		FilledAmount: big.NewInt(100),
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	err = dao.Create(o)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	updated := &types.Order{
		ID:           o.ID,
		UserAddress:  o.UserAddress,
		BaseToken:    o.BaseToken,
		QuoteToken:   o.QuoteToken,
		PricePoint:   big.NewInt(4000),
		Amount:       big.NewInt(4000),
		FilledAmount: big.NewInt(200),
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Signature:    o.Signature,
		Hash:         o.Hash,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}

	err = dao.Update(
		o.ID,
		updated,
	)

	if err != nil {
		t.Errorf("Could not updated order from hash %v", err)
	}

	queried, err := dao.GetByHash(o.Hash)
	if err != nil {
		t.Errorf("Could not get order by hash")
	}

	testutils.CompareOrder(t, queried, updated)
}

func TestOrderDao1(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Errorf("Could not drop previous order state")
	}

	o := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:  common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		BaseToken:    common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:   common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		Amount:       big.NewInt(1000),
		FilledAmount: big.NewInt(100),
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	err = dao.Create(o)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	o1, err := dao.GetByHash(common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"))
	if err != nil {
		t.Errorf("Could not get order by hash")
	}

	testutils.CompareOrder(t, o, o1)

	o2, err := dao.GetByUserAddress(common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"))
	if err != nil {
		t.Errorf("Could not get order by user address")
	}

	testutils.CompareOrder(t, o, o2[0])
}

func TestOrderDaoGetByHashes(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order state")
	}

	o1 := testutils.GetTestOrder1()
	o2 := testutils.GetTestOrder2()
	o3 := testutils.GetTestOrder3()

	dao.Create(&o1)
	dao.Create(&o2)
	dao.Create(&o3)

	orders, err := dao.GetByHashes([]common.Hash{o1.Hash, o2.Hash})
	if err != nil {
		t.Error("Could not get order by hashes")
	}

	assert.Equal(t, len(orders), 2)
	testutils.CompareOrder(t, orders[0], &o1)
	testutils.CompareOrder(t, orders[1], &o2)
}

func TestGetHistoryByUserAddress(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		FilledAmount: big.NewInt(0),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	o2 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0002"),
		UserAddress:  user,
		FilledAmount: big.NewInt(0),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	o3 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0003"),
		UserAddress:  user,
		FilledAmount: units.Ethers(5),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "PARTIAL_FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	o4 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0004"),
		UserAddress:  user,
		FilledAmount: units.Ethers(10),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		PricePoint:   big.NewInt(1e18),
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	dao.Create(o1)
	dao.Create(o2)
	dao.Create(o3)
	dao.Create(o4)

	lockedBalance, err := dao.GetUserLockedBalance(user, sellToken)
	if err != nil {
		t.Error("Could not get locked balance", err)
	}

	assert.Equal(t, units.Ethers(25), lockedBalance)
}

func TestGetUserOrderHistory(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		FilledAmount: big.NewInt(0),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	o2 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0002"),
		UserAddress:  user,
		FilledAmount: big.NewInt(0),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	o3 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0003"),
		UserAddress:  user,
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		FilledAmount: units.Ethers(5),
		Status:       "INVALID",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	o4 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0004"),
		UserAddress:  user,
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		FilledAmount: units.Ethers(10),
		PricePoint:   big.NewInt(1e18),
		Status:       "PARTIAL_FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         common.HexToHash("0x12"),
	}

	dao.Create(o1)
	dao.Create(o2)
	dao.Create(o3)
	dao.Create(o4)

	orders, err := dao.GetHistoryByUserAddress(user)
	if err != nil {
		t.Error("Could not get order history", err)
	}

	assert.Equal(t, 2, len(orders))
	testutils.CompareOrder(t, orders[0], o1)
	testutils.CompareOrder(t, orders[1], o3)
	assert.NotContains(t, orders, o2)
	assert.NotContains(t, orders, o4)
}

func TestUpdateOrderFilledAmount1(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")
	hash := common.HexToHash("0x5")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Amount:       units.Ethers(10),
		FilledAmount: big.NewInt(0),
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         hash,
	}

	err = dao.Create(o1)
	if err != nil {
		t.Error("Could not create order")
	}

	err = dao.UpdateOrderFilledAmount(hash, big.NewInt(5))
	if err != nil {
		t.Error("Could not get order history", err)
	}

	stored, err := dao.GetByHash(hash)
	if err != nil {
		t.Error("Could not retrieve order", err)
	}

	assert.Equal(t, "PARTIAL_FILLED", stored.Status)
	assert.Equal(t, big.NewInt(5), stored.FilledAmount)
}

func TestUpdateOrderFilledAmount2(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")
	hash := common.HexToHash("0x5")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		Amount:       units.Ethers(10),
		FilledAmount: units.Ethers(5),
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         hash,
	}

	err = dao.Create(o1)
	if err != nil {
		t.Error("Could not create order")
	}

	err = dao.UpdateOrderFilledAmount(hash, units.Ethers(6))
	if err != nil {
		t.Error("Could not get order history", err)
	}

	stored, err := dao.GetByHash(hash)
	if err != nil {
		t.Error("Could not retrieve order", err)
	}

	assert.Equal(t, "FILLED", stored.Status)
	assert.Equal(t, units.Ethers(10), stored.FilledAmount)
}

func TestUpdateOrderFilledAmount3(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")
	hash := common.HexToHash("0x5")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		Amount:       units.Ethers(10),
		FilledAmount: units.Ethers(5),
		Basetoken:    baseToken,
		QuoteToken:   quoteToken,
		Status:       "OPEN",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(50),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(50),
		Hash:         hash,
	}

	err = dao.Create(o1)
	if err != nil {
		t.Error("Could not create order")
	}

	err = dao.UpdateOrderFilledAmount(hash, math.Neg(units.Ethers(6)))
	if err != nil {
		t.Error("Could not get order history", err)
	}

	stored, err := dao.GetByHash(hash)
	if err != nil {
		t.Error("Could not retrieve order", err)
	}

	assert.Equal(t, "OPEN", stored.Status)
	assert.Equal(t, big.NewInt(0), stored.FilledAmount)
}

func TestUpdateOrderFilledAmounts(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")
	hash1 := common.HexToHash("0x5")
	hash2 := common.HexToHash("0x6")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		Amount:       units.Ethers(2),
		FilledAmount: units.Ethers(0),
		Status:       "FILLED",
		Side:         "BUY",
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(0),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(0),
		Hash:         hash1,
	}

	o2 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		Amount:       units.Ethers(1),
		FilledAmount: units.Ethers(0),
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(0),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(0),
		Hash:         hash2,
	}

	err = dao.Create(o1)
	if err != nil {
		t.Error("Could not create order")
	}

	err = dao.Create(o2)
	if err != nil {
		t.Error("Could not create order")
	}

	hashes := []common.Hash{hash1, hash2}
	amounts := []*big.Int{big.NewInt(-1), big.NewInt(-2)}
	orders, err := dao.UpdateOrderFilledAmounts(hashes, amounts)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(orders))
	assert.Equal(t, big.NewInt(1), orders[0].FilledAmount)
	assert.Equal(t, big.NewInt(2), orders[1].FilledAmount)
}

func TestOrderStatusesByHashes(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Error("Could not drop previous order collection")
	}

	user := common.HexToAddress("0x1")
	exchange := common.HexToAddress("0x2")
	baseToken := common.HexToAddress("0x3")
	quoteToken := common.HexToAddress("0x4")
	hash1 := common.HexToHash("0x5")
	hash2 := common.HexToHash("0x6")

	o1 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:  user,
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Amount:       units.Ethers(1),
		FilledAmount: units.Ethers(1),
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(0),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(0),
		Hash:         hash1,
	}

	o2 := &types.Order{
		ID:           bson.ObjectIdHex("537f700b537461b70c5f0002"),
		UserAddress:  user,
		BaseToken:    baseToken,
		QuoteToken:   quoteToken,
		Amount:       units.Ethers(1),
		FilledAmount: units.Ethers(1),
		Status:       "FILLED",
		Side:         "BUY",
		PairName:     "ZRX/WETH",
		MakeFee:      big.NewInt(0),
		Nonce:        big.NewInt(1000),
		TakeFee:      big.NewInt(0),
		Hash:         hash2,
	}

	err = dao.Create(o1)
	if err != nil {
		t.Error("Could not create order")
	}

	err = dao.Create(o2)
	if err != nil {
		t.Error("Could not create order")
	}

	orders, err := dao.UpdateOrderStatusesByHashes("INVALIDATED", hash1, hash2)
	if err != nil {
		t.Error("Error in updateOrderStatusHashes", err)
	}

	assert.Equal(t, 2, len(orders))
	assert.Equal(t, "INVALIDATED", orders[0].Status)
	assert.Equal(t, "INVALIDATED", orders[1].Status)
}

func ExampleGetOrderBook() {
	session, err := mgo.Dial(app.Config.DSN)
	if err != nil {
		panic(err)
	}

	db = &Database{session}
	pairDao := NewPairDao(PairDaoDBOption("proofdex"))
	orderDao := NewOrderDao(OrderDaoDBOption("proofdex"))
	pair, err := pairDao.GetByTokenSymbols("BAT", "WETH")
	if err != nil {
		panic(err)
	}

	bids, asks, err := orderDao.GetOrderBook(pair)
	if err != nil {
		panic(err)
	}

	utils.PrintJSON(bids)
	utils.PrintJSON(asks)
}

func ExampleGetOrderBookPricePoint() {
	session, err := mgo.Dial(app.Config.DSN)
	if err != nil {
		panic(err)
	}

	db = &Database{session}

	pairDao := NewPairDao(PairDaoDBOption("proofdex"))
	orderDao := NewOrderDao(OrderDaoDBOption("proofdex"))
	pair, err := pairDao.GetByTokenSymbols("AE", "WETH")
	if err != nil {
		panic(err)
	}

	orderPricePoint, err := orderDao.GetOrderBookPricePoint(pair, big.NewInt(59303), "BUY")
	if err != nil {
		panic(err)
	}

	utils.PrintJSON(orderPricePoint)
}

func ExampleGetRawOrderBook() {
	session, err := mgo.Dial(app.Config.DSN)
	if err != nil {
		panic(err)
	}

	db = &Database{session}

	pairDao := NewPairDao(PairDaoDBOption("proofdex"))
	orderDao := NewOrderDao(OrderDaoDBOption("proofdex"))
	pair, err := pairDao.GetByTokenSymbols("AE", "WETH")
	if err != nil {
		panic(err)
	}

	orders, err := orderDao.GetRawOrderBook(pair)
	if err != nil {
		panic(err)
	}

	utils.PrintJSON(orders)
}
