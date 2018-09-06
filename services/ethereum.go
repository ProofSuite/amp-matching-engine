package services

// import (
// 	"context"
// 	"log"
// 	"math/big"

// 	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
// 	"github.com/Proofsuite/amp-matching-engine/interfaces"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethTypes "github.com/ethereum/go-ethereum/core/types"
// )

// type EthereumService struct {
// 	EthereumClient interfaces.EthereumClient
// }

// func NewEthereumService(e interfaces.EthereumClient) *EthereumService {
// 	return &EthereumService{e}
// }

// func (s *EthereumService) WaitMined(tx *ethTypes.Transaction) (*ethTypes.Receipt, error) {
// 	ctx := context.Background()
// 	receipt, err := bind.WaitMined(ctx, s.EthereumClient, tx)
// 	if err != nil {
// 		log.Print(err)
// 		return &ethTypes.Receipt{}, err
// 	}

// 	return receipt, nil
// }

// func (s *EthereumService) GetBalanceAt(a common.Address) (*big.Int, error) {
// 	ctx := context.Background()
// 	nonce, err := s.EthereumClient.BalanceAt(ctx, a, nil)
// 	if err != nil {
// 		log.Print(err)
// 		return big.NewInt(0), err
// 	}

// 	return nonce, nil
// }

// func (s *EthereumService) GetPendingNonceAt(a common.Address) (uint64, error) {
// 	ctx := context.Background()
// 	nonce, err := s.EthereumClient.PendingNonceAt(ctx, a)
// 	if err != nil {
// 		log.Print(err)
// 		return 0, err
// 	}

// 	return nonce, nil
// }

// func (s *EthereumService) BalanceOf(owner common.Address, token common.Address) (*big.Int, error) {
// 	tokenInterface, err := contractsinterfaces.NewToken(token, s.EthereumClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	opts := &bind.CallOpts{Pending: true}
// 	b, err := tokenInterface.BalanceOf(opts, owner)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	return b, nil
// }

// func (s *EthereumService) Allowance(owner, spender, token common.Address) (*big.Int, error) {
// 	tokenInterface, err := contractsinterfaces.NewToken(token, s.EthereumClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	opts := &bind.CallOpts{Pending: true}
// 	a, err := tokenInterface.Allowance(opts, owner, spender)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	return a, nil
// }
