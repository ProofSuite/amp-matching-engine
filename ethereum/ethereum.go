package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Conn is singleton rabbitmq connection
var ethereumClient *ethclient.Client

// InitConnection Initializes single rabbitmq connection for whole system
func InitConnection(url string) {

	rpcClient, err := rpc.DialHTTP(url)
	if err != nil {
		panic(err)
	}

	ethereumClient = ethclient.NewClient(rpcClient)
}
