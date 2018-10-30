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
}

func NewValidatorService(
	ethereumProvider interfaces.EthereumProvider,
	accountDao interfaces.AccountDao,
	orderDao interfaces.OrderDao,
) *ValidatorService {

	return &ValidatorService{
		ethereumProvider,
		accountDao,
		orderDao,
	}
}

func (s *ValidatorService) ValidateBalance(o *types.Order) error {
	wethAddress := common.HexToAddress(app.Config.Ethereum["weth_address"])
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	balanceRecord, err := s.accountDao.GetTokenBalances(o.UserAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	wethBalance, err := s.ethereumProvider.BalanceOf(o.UserAddress, wethAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	wethAllowance, err := s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, wethAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenBalance, err := s.ethereumProvider.BalanceOf(o.UserAddress, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenAllowance, err := s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	wethLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, wethAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	fee := math.Max(o.MakeFee, o.TakeFee)
	availableWethBalance := math.Sub(wethBalance, wethLockedBalance)
	availableSellTokenBalance := math.Sub(sellTokenBalance, sellTokenLockedBalance)

	//WETH Token Balance (for fees)
	if availableWethBalance.Cmp(fee) == -1 {
		return errors.New("Insufficient WETH Balance")
	}

	if wethAllowance.Cmp(fee) == -1 {
		return errors.New("Insufficient WETH Allowance")
	}

	//Sell Token Balance
	if sellTokenBalance.Cmp(o.SellAmount) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	if availableSellTokenBalance.Cmp(o.SellAmount) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	if sellTokenAllowance.Cmp(o.SellAmount) == -1 {
		return fmt.Errorf("Insufficient %v Allowance", o.SellTokenSymbol())
	}

	utils.PrintJSON(sellTokenAllowance)

	sellTokenBalanceRecord := balanceRecord[o.SellToken]
	if sellTokenBalanceRecord == nil {
		return errors.New("Account error: Balance record not found")
	}

	wethTokenBalanceRecord := balanceRecord[wethAddress]
	if wethTokenBalanceRecord == nil {
		return errors.New("Account error: Balance record not found")
	}

	sellTokenBalanceRecord.Balance.Set(sellTokenBalance)
	sellTokenBalanceRecord.Allowance.Set(sellTokenAllowance)
	wethTokenBalanceRecord.Balance.Set(wethBalance)
	wethTokenBalanceRecord.Allowance.Set(wethAllowance)

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, wethAddress, wethTokenBalanceRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken, sellTokenBalanceRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
