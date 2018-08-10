package ethereumclient

import (
	"context"

	"github.com/Proofsuite/go-ethereum/ethclient"
	"github.com/Proofsuite/go-ethereum/rpc"
)

// Conn is singleton rabbitmq connection
var EthereumClient *ethclient.Client

// InitConnection Initializes single rabbitmq connection for whole system
func InitConnection(url string) {

	rpcClient, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)
}
