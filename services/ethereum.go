package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService struct {
	EthereumClient *ethclient.Client
}

type EthereumServiceInterface interface {
	WaitMined(tx *ethTypes.Transaction) (*ethTypes.Receipt, error)
	GetPendingBalanceAt(a common.Address) (*big.Int, error)
}

func NewEthereumService(e *ethclient.Client) *EthereumService {
	return &EthereumService{e}
}

func (s *EthereumService) WaitMined(tx *ethTypes.Transaction) (*ethTypes.Receipt, error) {
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, s.EthereumClient, tx)

	if err != nil {
		return &ethTypes.Receipt{}, err
	}

	return receipt, nil
}

func (s *EthereumService) GetPendingBalanceAt(a common.Address) (*big.Int, error) {
	ctx := context.Background()
	balance, err := s.EthereumClient.PendingBalanceAt(ctx, a)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
