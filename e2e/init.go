package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/crons"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/endpoints"
	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/operator"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/dbtest"
)

type apiTestCase struct {
	tag         string
	method      string
	url         string
	body        string
	status      int
	response    interface{}
	checkMethod string
	compareFn   func(t *testing.T, actual, expected interface{})
}

var server dbtest.DBServer

// Init function initializes the e2e testing
func Init(t *testing.T) {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	_, err := daos.InitSession(session)
	if err != nil {
		panic(err)
	} else {
		err = session.DB(app.Config.DBName).DropDatabase()
	}

	// === drop database on test end ===
	defer session.DB(app.Config.DBName).DropDatabase()
	tokens := testToken(t)
	pair := testPair(t, tokens)
	accounts := testAccount(t, tokens)
	testWS(t, pair, accounts)
	// address := testAddress(t, tokens)
	// testBalance(t, tokens, address)
}

func NewRouter() *routing.Router {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	router := routing.New()

	router.To("GET,HEAD", "/ping", func(c *routing.Context) error {
		c.Abort() // skip all other middlewares/handlers
		return c.Write("OK " + app.Version)
	})

	router.Use(
		app.Init(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)

	rg := router.Group("")

	provider := ethereum.NewWebsocketProvider()
	amqp := rabbitmq.InitConnection(app.Config.Rabbitmq)
	redisClient := redis.NewRedisConnection(app.Config.Redis)
	redisClient.FlushAll()

	eng, err := engine.InitEngine(redisClient, amqp)
	if err != nil {
		panic(err)
	}

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	tradeDao := daos.NewTradeDao()
	accountDao := daos.NewAccountDao()
	walletDao := daos.NewWalletDao()

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao)
	ohlcvService := services.NewOHLCVService(tradeDao)
	tokenService := services.NewTokenService(tokenDao)
	tradeService := services.NewTradeService(tradeDao)
	pairService := services.NewPairService(pairDao, tokenDao, eng, tradeService)
	orderService := services.NewOrderService(orderDao, pairDao, accountDao, tradeDao, eng, provider, amqp)
	orderBookService := services.NewOrderBookService(pairDao, tokenDao, eng)
	walletService := services.NewWalletService(walletDao)
	cronService := crons.NewCronService(ohlcvService)

	// get exchange contract instance
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	exchange, err := contracts.NewExchange(
		walletService,
		exchangeAddress,
		provider.Client,
	)

	if err != nil {
		panic(err)
	}

	// deploy operator
	op, err := operator.NewOperator(
		walletService,
		tradeService,
		orderService,
		provider,
		exchange,
	)

	if err != nil {
		panic(err)
	}

	endpoints.ServeAccountResource(rg, accountService)
	endpoints.ServeTokenResource(rg, tokenService)
	endpoints.ServePairResource(rg, pairService)
	endpoints.ServeOrderBookResource(rg, orderBookService)
	endpoints.ServeOHLCVResource(rg, ohlcvService)
	endpoints.ServeTradeResource(rg, tradeService)
	endpoints.ServeOrderResource(rg, orderService, eng)

	//initialize rabbitmq subscriptions
	orderService.SubscribeOrders(eng.HandleOrders)
	orderService.SubscribeTrades(op.HandleTrades)
	eng.SubscribeResponseQueue(orderService.HandleEngineResponse)

	// fmt.Printf("\n%+v\n", app.Config.TickDuration)
	cronService.InitCrons()
	return router
}

func testAPI(router *routing.Router, method, URL, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, URL, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	httptest.NewServer(router)
	return res
}

func runAPITests(t *testing.T, router *routing.Router, tests []apiTestCase) {
	for _, test := range tests {
		res := testAPI(router, test.method, test.url, test.body)
		if test.response != "" {
			var resp interface{}
			if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
				fmt.Printf("%v", err)
			}
			switch test.checkMethod {
			case "contains":
				assert.Contains(t, resp, test.response, test.tag)
			case "equals":
				assert.JSONEq(t, test.response.(string), res.Body.String(), test.tag)
			case "subset":
				assert.Subset(t, resp, test.response, test.tag)
			case "custom":
				test.compareFn(t, resp, test.response)
			}
		}
	}
}
