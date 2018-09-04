package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/crons"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/endpoints"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Proofsuite/amp-matching-engine/engine"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Get application up and running",
	Long:  `Get application up and running`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func run(cmd *cobra.Command, args []string) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")
	logger := logrus.New()

	// connect to the database
	if _, err := daos.InitSession(nil); err != nil {
		panic(err)
	}

	http.Handle("/", NewRouter(logger))
	http.HandleFunc("/socket", ws.ConnectionEndpoint)

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	logger.Infof("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func NewRouter(logger *logrus.Logger) *routing.Router {
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

	rabbitmq.InitConnection(app.Config.Rabbitmq)
	ethereum.InitConnection(app.Config.Ethereum)
	redisClient := redis.NewRedisConnection(app.Config.Redis)

	// instantiate engine
	eng, err := engine.InitEngine(redisClient)
	if err != nil {
		panic(err)
	}

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	tradeDao := daos.NewTradeDao()
	accountDao := daos.NewAccountDao()

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao)
	ohlcvService := services.NewOHLCVService(tradeDao)
	tokenService := services.NewTokenService(tokenDao)
	tradeService := services.NewTradeService(tradeDao)
	pairService := services.NewPairService(pairDao, tokenDao, eng, tradeService)
	orderService := services.NewOrderService(orderDao, pairDao, accountDao, tradeDao, eng)
	orderBookService := services.NewOrderBookService(pairDao, tokenDao, eng)
	cronService := crons.NewCronService(ohlcvService)
	// walletService := services.NewWalletService(walletDao, balanceDao)

	endpoints.ServeAccountResource(rg, accountService)
	endpoints.ServeTokenResource(rg, tokenService)
	endpoints.ServePairResource(rg, pairService)
	endpoints.ServeOrderBookResource(rg, orderBookService)
	endpoints.ServeOHLCVResource(rg, ohlcvService)
	endpoints.ServeTradeResource(rg, tradeService)
	endpoints.ServeOrderResource(rg, orderService, eng)

	//initialize rabbitmq subscriptions
	orderService.SubscribeQueue(eng.HandleOrders)
	eng.SubscribeResponseQueue(orderService.HandleEngineResponse)

	// fmt.Printf("\n%+v\n", app.Config.TickDuration)
	cronService.InitCrons()
	return router
}
