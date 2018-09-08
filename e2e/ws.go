package e2e

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func testWS(t *testing.T, pairs []types.Pair, accounts map[*ecdsa.PrivateKey]types.Account) {
	fmt.Printf("\n=== Starting WS tests ===\n")
	exchangeAddress := common.HexToAddress("0x")
	pair := &pairs[0]
	baseToken := pair.BaseTokenAddress
	quoteToken := pair.QuoteTokenAddress
	wallets := make([]*types.Wallet, 0)
	clients := make([]*testutils.Client, 0)
	factories := make([]*testutils.OrderFactory, 0)

	for key := range accounts {
		w := types.NewWalletFromPrivateKey(hex.EncodeToString(crypto.FromECDSA(key)))
		c := testutils.NewClient(w, http.HandlerFunc(ws.ConnectionEndpoint))
		f, err := testutils.NewOrderFactory(pair, w, exchangeAddress)
		if err != nil {
			panic(err)
		}
		wallets = append(wallets, w)
		clients = append(clients, c)
		factories = append(factories, f)
		c.Start()
	}

	obClient := newObClient(t, baseToken, quoteToken, nil)
	tradeClient := newTradeClient(t, baseToken, quoteToken, nil)

	// check if ohlcv client gets UPDATE payload in 5secs
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(test *testing.T) {
		ohlcvClient := newOHLCVClient(t, baseToken, quoteToken, nil)
		time.Sleep(6 * time.Second)
		log := getLatestRLog(ohlcvClient.ResponseLogs)
		assert.Equal(test, "UPDATE", log.Payload.Type)
		wg.Done()
	}(t)

	clients = append(clients, obClient)
	clients = append(clients, tradeClient)

	NewRouter()
	testInitSubscription(t, clients[0], factories[0], clients[1], factories[1], tradeClient, baseToken, quoteToken)

	wg.Wait()
}

func testInitSubscription(t *testing.T, client1 *testutils.Client, factory1 *testutils.OrderFactory, client2 *testutils.Client, factory2 *testutils.OrderFactory, tradeClient *testutils.Client, baseToken, quoteToken common.Address) {
	fmt.Printf("\t== testInitSubscription ==\n")
	// send buy order
	buyOrderMsg, _, err := factory1.NewBuyOrderMessage(1e5, 1e6)
	if err != nil {
		panic(err)
	}
	client1.Requests <- buyOrderMsg
	time.Sleep(time.Second)
	assert.Equal(t, "ORDER_ADDED", getLatestRLog(client1.ResponseLogs).Payload.Type)

	// send sell order
	sellOrderMsg, _, err := factory2.NewSellOrderMessage(1e5+10, 1e6)
	if err != nil {
		panic(err)
	}
	client2.Requests <- sellOrderMsg
	time.Sleep(time.Second)
	assert.Equal(t, "ORDER_ADDED", getLatestRLog(client2.ResponseLogs).Payload.Type)

	newTradeClient(t, baseToken, quoteToken, getLatestRLog(tradeClient.ResponseLogs).Payload.Data)
}

func getOrderbookSubscribeRequest(baseToken, quoteToken common.Address) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Channel: ws.LiteOrderBookChannel,
		Payload: types.WebSocketPayload{
			Type: "subscription",
			Data: types.WebSocketSubscription{
				Event: types.SUBSCRIBE,
				Pair: types.PairSubDoc{
					BaseToken:  baseToken,
					QuoteToken: quoteToken,
				},
			},
		},
	}
}

func getTradeSubscribeRequest(baseToken, quoteToken common.Address) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Channel: ws.TradeChannel,
		Payload: types.WebSocketPayload{
			Type: "subscription",
			Data: types.WebSocketSubscription{
				Event: types.SUBSCRIBE,
				Pair: types.PairSubDoc{
					BaseToken:  baseToken,
					QuoteToken: quoteToken,
				},
			},
		},
	}
}
func getOHLCVSubscribeRequest(baseToken, quoteToken common.Address) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Channel: ws.OHLCVChannel,
		Payload: types.WebSocketPayload{
			Type: "subscription",
			Data: types.WebSocketSubscription{
				Event: types.SUBSCRIBE,
				Pair: types.PairSubDoc{
					BaseToken:  baseToken,
					QuoteToken: quoteToken,
				},
				Params: types.Params{
					Duration: 5,
					Units:    "sec",
				},
			},
		},
	}
}
func getWebsocketMessage(channel, t, hash string, data interface{}) types.WebSocketMessage {
	return types.WebSocketMessage{
		Channel: channel,
		Payload: types.WebSocketPayload{
			Type: t,
			Hash: "",
			Data: data,
		},
	}
}

func getLatestRLog(logs []types.WebSocketMessage) types.WebSocketMessage {
	return logs[len(logs)-1]
}

func resetLogs(clients ...*testutils.Client) {
	for _, client := range clients {
		client.ResponseLogs = make([]types.WebSocketMessage, 0)
		client.RequestLogs = make([]types.WebSocketMessage, 0)
	}
}

func newObClient(t *testing.T, baseToken, quoteToken common.Address, testData interface{}) *testutils.Client {
	// orderBook client
	k, _ := crypto.GenerateKey()
	w := types.NewWalletFromPrivateKey(hex.EncodeToString(crypto.FromECDSA(k)))
	obClient := testutils.NewClient(w, http.HandlerFunc(ws.ConnectionEndpoint))
	obClient.Start()

	// Subscribe to orderbook channel
	obClient.Requests <- getOrderbookSubscribeRequest(baseToken, quoteToken)
	time.Sleep(time.Second)

	if testData == nil {
		testData = map[string]interface{}{
			"asks": nil,
			"bids": nil,
		}
	}

	expectedRes := getWebsocketMessage(ws.LiteOrderBookChannel, "INIT", "", testData)

	assert.Equal(t, expectedRes, obClient.ResponseLogs[0])
	return obClient
}
func newTradeClient(t *testing.T, baseToken, quoteToken common.Address, testData interface{}) *testutils.Client {
	//tradeClient
	k, _ := crypto.GenerateKey()
	w := types.NewWalletFromPrivateKey(hex.EncodeToString(crypto.FromECDSA(k)))
	tradeClient := testutils.NewClient(w, http.HandlerFunc(ws.ConnectionEndpoint))
	tradeClient.Start()

	// Subscribe to trades channel
	tradeClient.Requests <- getTradeSubscribeRequest(baseToken, quoteToken)
	time.Sleep(time.Second)

	expectedRes := getWebsocketMessage(ws.TradeChannel, "INIT", "", testData)
	assert.Equal(t, expectedRes, tradeClient.ResponseLogs[0])

	return tradeClient
}

func newOHLCVClient(t *testing.T, baseToken, quoteToken common.Address, testData interface{}) *testutils.Client {
	//ohlcvClient
	k, _ := crypto.GenerateKey()
	w := types.NewWalletFromPrivateKey(hex.EncodeToString(crypto.FromECDSA(k)))
	ohlcvClient := testutils.NewClient(w, http.HandlerFunc(ws.ConnectionEndpoint))
	ohlcvClient.Start()

	// Subscribe to trades channel
	ohlcvClient.Requests <- getOHLCVSubscribeRequest(baseToken, quoteToken)
	time.Sleep(time.Second)

	expectedRes := getWebsocketMessage(ws.OHLCVChannel, "INIT", "", testData)
	assert.Equal(t, expectedRes, ohlcvClient.ResponseLogs[0])
	return ohlcvClient
}

func signTrades(trades []*types.Trade, wallet *types.Wallet) {
	for _, trade := range trades {
		if err := wallet.SignTrade(trade); err != nil {
			panic(err)
		}
	}
}

func print(data interface{}, tag string) {
	dab, _ := json.Marshal(data)
	fmt.Printf("\n XXX %s XXX \n>>>> %s\n", tag, string(dab))
}
