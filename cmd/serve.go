package cmd

import (
	"fmt"
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/crons"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/endpoints"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/operator"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/Proofsuite/go-ethereum/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/Proofsuite/amp-matching-engine/engine"
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
	// connect to the database
	_, err := daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	rabbitConn := rabbitmq.InitConnection(app.Config.Rabbitmq)
	redisConn := redis.NewRedisConnection(app.Config.Redis)
	provider := ethereum.NewWebsocketProvider()

	router := NewRouter(provider, redisConn, rabbitConn)
	http.Handle("/", router)
	http.HandleFunc("/socket", ws.ConnectionEndpoint)

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Info("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func NewRouter(
	provider *ethereum.EthereumProvider,
	redisConn *redis.RedisConnection,
	rabbitConn *rabbitmq.Connection,
) *mux.Router {

	r := mux.NewRouter()

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	tradeDao := daos.NewTradeDao()
	accountDao := daos.NewAccountDao()
	walletDao := daos.NewWalletDao()

	// instantiate engine
	eng := engine.NewEngine(redisConn, rabbitConn, pairDao)

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao)
	ohlcvService := services.NewOHLCVService(tradeDao)
	tokenService := services.NewTokenService(tokenDao)
	tradeService := services.NewTradeService(tradeDao)
	pairService := services.NewPairService(pairDao, tokenDao, eng, tradeService)
	orderService := services.NewOrderService(orderDao, pairDao, accountDao, tradeDao, eng, provider, rabbitConn)
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
		rabbitConn,
	)

	if err != nil {
		panic(err)
	}

	// deploy http and ws endpoints
	endpoints.ServeAccountResource(r, accountService)
	endpoints.ServeTokenResource(r, tokenService)
	endpoints.ServePairResource(r, pairService)
	endpoints.ServeOrderBookResource(r, orderBookService)
	endpoints.ServeOHLCVResource(r, ohlcvService)
	endpoints.ServeTradeResource(r, tradeService)
	endpoints.ServeOrderResource(r, orderService, eng)

	//initialize rabbitmq subscriptions
	rabbitConn.SubscribeOrders(eng.HandleOrders)
	rabbitConn.SubscribeTrades(op.HandleTrades)
	rabbitConn.SubscribeOperator(orderService.HandleOperatorMessages)
	rabbitConn.SubscribeEngineResponses(orderService.HandleEngineResponse)

	cronService.InitCrons()
	return r
}
