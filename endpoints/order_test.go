package endpoints

import (
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redisclient"
	"github.com/Proofsuite/amp-matching-engine/services"
)

func TestOrder(t *testing.T) {

	rabbitmq.InitConnection("amqp://guest:guest@localhost:5672/")

	router := newRouter()

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	pairDao := daos.NewPairDao()
	balanceDao := daos.NewBalanceDao()
	tradesDao := daos.NewTradeDao()

	// instantiate engine
	e, err := engine.InitEngine(orderDao, redisclient.InitConnection(app.Config.Redis))
	if err != nil {
		panic(err)
	}

	// get services for injection
	orderService := services.NewOrderService(orderDao, balanceDao, pairDao, tradesDao, e)

	ServeOrderResource(&router.RouteGroup, orderService, e)

	a := map[string]interface{}{
		"channel": "order_book",
		"message": map[string]string{
			"event": "subscribe",
			"key":   "hpc-aut",
		},
	}
	testSocket(a)
}
