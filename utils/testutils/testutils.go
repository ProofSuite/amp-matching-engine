package testutils

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
)

func Mine(client *ethereum.SimulatedClient) {
	nextTime := time.Now()
	nextTime = nextTime.Add(500 * time.Millisecond)
	time.Sleep(time.Until(nextTime))

	client.Commit()
	go Mine(client)
}

type MockServices struct {
	WalletService    *mocks.WalletService
	AccountService   *mocks.AccountService
	EthereumService  *mocks.EthereumService
	OrderService     *mocks.OrderService
	OrderBookService *mocks.OrderBookService
	TokenService     *mocks.TokenService
	TxService        *mocks.TxService
	PairService      *mocks.PairService
	TradeService     *mocks.TradeService
}

type MockDaos struct {
	WalletDao  *mocks.WalletDao
	AccountDao *mocks.AccountDao
	OrderDao   *mocks.OrderDao
	TokenDao   *mocks.TokenDao
	TradeDao   *mocks.TradeDao
	PairDao    *mocks.PairDao
}

func NewMockServices() *MockServices {
	return &MockServices{
		WalletService:    new(mocks.WalletService),
		AccountService:   new(mocks.AccountService),
		EthereumService:  new(mocks.EthereumService),
		OrderService:     new(mocks.OrderService),
		OrderBookService: new(mocks.OrderBookService),
		TokenService:     new(mocks.TokenService),
		TxService:        new(mocks.TxService),
		PairService:      new(mocks.PairService),
	}
}

func NewMockDaos() *MockDaos {
	return &MockDaos{
		WalletDao:  new(mocks.WalletDao),
		AccountDao: new(mocks.AccountDao),
		OrderDao:   new(mocks.OrderDao),
		TokenDao:   new(mocks.TokenDao),
		TradeDao:   new(mocks.TradeDao),
		PairDao:    new(mocks.PairDao),
	}
}
