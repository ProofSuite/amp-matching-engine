package services

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type EthereumService struct {
	EthereumClient EthereumClientInterface
}

type EthereumClientInterface interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	// PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
}

func NewEthereumService(e EthereumClientInterface) *EthereumService {
	return &EthereumService{e}
}

func (s *EthereumService) WaitMined(tx *ethTypes.Transaction) (*ethTypes.Receipt, error) {
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, s.EthereumClient, tx)
	if err != nil {
		log.Print(err)
		return &ethTypes.Receipt{}, err
	}

	return receipt, nil
}

// func (s *EthereumService) GetPendingBalanceAt(a common.Address) (*big.Int, error) {
// 	ctx := context.Background()
// 	balance, err := s.EthereumClient.PendingBalanceAt(ctx, a)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return balance, nil
// }

func (s *EthereumService) GetPendingNonceAt(a common.Address) (uint64, error) {
	ctx := context.Background()
	nonce, err := s.EthereumClient.PendingNonceAt(ctx, a)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return nonce, nil
}
