package main

import (
	"fmt"
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/endpoints"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redisclient"
	"github.com/Proofsuite/amp-matching-engine/services"

	"github.com/Proofsuite/amp-matching-engine/engine"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Sirupsen/logrus"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
)

func main() {
	// load application configurations
	if err := app.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	// load error messages
	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	// create the logger
	logger := logrus.New()

	rabbitmq.InitConnection(app.Config.Rabbitmq)

	// connect to the database
	if err := daos.InitSession(); err != nil {
		panic(err)
	}

	// wire up API routing
	http.Handle("/", buildRouter(logger))

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	logger.Infof("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func buildRouter(logger *logrus.Logger) *routing.Router {
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

	// rg.Post("/auth", apis.Auth(app.Config.JWTSigningKey))
	// rg.Use(auth.JWT(app.Config.JWTVerificationKey, auth.JWTOptions{
	// 	SigningMethod: app.Config.JWTSigningMethod,
	// 	TokenHandler:  apis.JWTHandler,
	// }))

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
	orderService := services.NewOrderService(orderDao, balanceDao, pairDao, tradesDao)
	addressService := services.NewAddressService(addressDao, balanceDao, tokenDao)

	endpoints.ServeTokenResource(rg, tokenService)
	endpoints.ServePairResource(rg, pairService)
	endpoints.ServeBalanceResource(rg, balanceService)
	endpoints.ServeOrderResource(rg, orderService, e)
	endpoints.ServeTradeResource(rg, tradeService)
	endpoints.ServeAddressResource(rg, addressService)

	return router
}
