package services

import (
	"errors"
	"fmt"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
)

type ValidatorService struct {
	ethereumProvider interfaces.EthereumProvider
	accountDao       interfaces.AccountDao
	orderDao         interfaces.OrderDao
	pairDao          interfaces.PairDao
}

func NewValidatorService(
	ethereumProvider interfaces.EthereumProvider,
	accountDao interfaces.AccountDao,
	orderDao interfaces.OrderDao,
	pairDao interfaces.PairDao,
) *ValidatorService {

	return &ValidatorService{
		ethereumProvider,
		accountDao,
		orderDao,
		pairDao,
	}
}

func (s *ValidatorService) ValidateBalance(o *types.Order) error {
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	pricepointMultiplier := pair.PricepointMultiplier()

	utils.PrintJSON(pricepointMultiplier)
	utils.PrintJSON(o.SellAmount(pricepointMultiplier))

	balanceRecord, err := s.accountDao.GetTokenBalances(o.UserAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenBalance, err := s.ethereumProvider.BalanceOf(o.UserAddress, o.SellToken())
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenAllowance, err := s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, o.SellToken())
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, o.SellToken())
	if err != nil {
		logger.Error(err)
		return err
	}

	availableSellTokenBalance := math.Sub(sellTokenBalance, sellTokenLockedBalance)

	//Sell Token Balance
	if sellTokenBalance.Cmp(o.SellAmount(pricepointMultiplier)) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	if availableSellTokenBalance.Cmp(o.SellAmount(pricepointMultiplier)) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	if sellTokenAllowance.Cmp(o.SellAmount(pricepointMultiplier)) == -1 {
		return fmt.Errorf("Insufficient %v Allowance", o.SellTokenSymbol())
	}

	sellTokenBalanceRecord := balanceRecord[o.SellToken()]
	if sellTokenBalanceRecord == nil {
		return errors.New("Account error: Balance record not found")
	}

	sellTokenBalanceRecord.Balance.Set(sellTokenBalance)
	sellTokenBalanceRecord.Allowance.Set(sellTokenAllowance)

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken(), sellTokenBalanceRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
