package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/crons"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/endpoints"
	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/matching-engine/redisclient"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/stretchr/testify/assert"
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

// Init function initializes the e2e testing
func Init(t *testing.T) {
	// the test may be started from the home dire
	if session, err := daos.InitSession(); err != nil {
		panic(err)
	} else {
		err = session.DB(app.Config.DBName).DropDatabase()

		// === drop database on test end ===
		// defer session.DB(app.Config.DBName).DropDatabase()
	}

	tokens := testToken(t)
	testPair(t, tokens)
	address := testAddress(t, tokens)
	testBalance(t, tokens, address)

}

func buildRouter() *routing.Router {

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

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	balanceDao := daos.NewBalanceDao()
	addressDao := daos.NewAddressDao()
	tradesDao := daos.NewTradeDao()

	// instantiate engine
	e, err := engine.InitEngine(orderDao, redisclient.InitConnection(app.Config.Redis))
	if err != nil {
		panic(err)
	}

	// get services for injection
	tokenService := services.NewTokenService(tokenDao)
	tradeService := services.NewTradeService(tradesDao)
	pairService := services.NewPairService(pairDao, tokenDao, e, tradeService)
	balanceService := services.NewBalanceService(balanceDao, tokenDao)
	orderService := services.NewOrderService(orderDao, balanceDao, pairDao, tradesDao, addressDao, e)
	addressService := services.NewAddressService(addressDao, balanceDao, tokenDao)

	endpoints.ServeTokenResource(rg, tokenService)
	endpoints.ServePairResource(rg, pairService)
	endpoints.ServeBalanceResource(rg, balanceService)
	endpoints.ServeOrderResource(rg, orderService, e)
	endpoints.ServeTradeResource(rg, tradeService)
	endpoints.ServeAddressResource(rg, addressService)

	cronService := crons.NewCronService(tradeService)
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
