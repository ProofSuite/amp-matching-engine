package contracts_test

import (
	"log"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
)

func SetupTest() (*testutils.Deployer, *types.Wallet, common.Address, common.Address, *types.Wallet, *types.Wallet) {
	err := app.LoadConfig("../config")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	wallet := testutils.GetTestWallet1()
	maker := testutils.GetTestWallet2()
	taker := testutils.GetTestWallet3()

	walletDao := new(mocks.WalletDao)
	walletDao.On("GetDefaultAdminWallet").Return(wallet, nil)

	walletService := services.NewWalletService(walletDao)
	txService := services.NewTxService(walletDao, wallet)

	deployer, err := testutils.NewSimulator(walletService, txService, []common.Address{wallet.Address, maker.Address, taker.Address})
	if err != nil {
		panic(err)
	}

	feeAccount := common.HexToAddress(app.Config.FeeAccount)
	wethToken := common.HexToAddress(app.Config.WETH)

	return deployer, wallet, feeAccount, wethToken, maker, taker
}

func TestSetFeeAccount(t *testing.T) {
	deployer, _, feeAccount, wethToken, _, _ := SetupTest()
	exchange, _, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	newFeeAccount := testutils.GetTestAddress1()

	_, err = exchange.SetFeeAccount(newFeeAccount)
	if err != nil {
		t.Errorf("Could not see new fee account: %v", err)
	}

	simulator.Commit()

	feeAccount, err = exchange.FeeAccount()
	if err != nil {
		t.Errorf("Error retrieving fee account address: %v", err)
	}

	if newFeeAccount != feeAccount {
		t.Errorf("Fee account not set correctly")
	}
}

func TestSetOperator(t *testing.T) {
	deployer, _, feeAccount, wethToken, _, _ := SetupTest()

	exchange, _, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange")
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	operator := testutils.GetTestAddress1()

	_, err = exchange.SetOperator(operator, true)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	simulator.Commit()

	isOperator, err := exchange.Operator(operator)
	if err != nil {
		t.Errorf("Error calling the operator variable: %v", err)
	}

	if isOperator != true {
		t.Errorf("Operator variable should be equal to true but got false")
	}
}

func TestTrade(t *testing.T) {
	deployer, admin, feeAccount, wethToken, _, _ := SetupTest()

	maker := testutils.GetTestWallet1()
	taker := testutils.GetTestWallet2()
	buyAmount := big.NewInt(1e18)
	sellAmount := big.NewInt(1e18)
	amount := big.NewInt(5e17)
	expires := big.NewInt(1e7)

	exchange, exchangeAddr, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange")
	}

	_, err = exchange.SetOperator(admin.Address, true)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	//Initially Maker owns 1e18 units of sellToken and Taker owns 1e18 units buyToken
	sellToken, sellTokenAddr, _, err := deployer.DeployToken(maker.Address, sellAmount)
	if err != nil {
		t.Errorf("Error deploying token 1: %v", err)
	}

	buyToken, buyTokenAddr, _, err := deployer.DeployToken(taker.Address, buyAmount)
	if err != nil {
		t.Errorf("Error deploying token 2: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	exchange.PrintErrors()

	sellToken.SetTxSender(maker)
	_, err = sellToken.Approve(exchangeAddr, sellAmount)
	if err != nil {
		t.Errorf("Could not approve sellToken: %v", err)
	}

	buyToken.SetTxSender(taker)
	_, err = buyToken.Approve(exchangeAddr, buyAmount)
	if err != nil {
		t.Errorf("Could not approve buyToken: %v", err)
	}

	simulator.Commit()

	//Maker creates an order that exchanges 'sellAmount' of sellToken for 'buyAmount' of buyToken
	order := &types.Order{
		ExchangeAddress: exchangeAddr,
		BuyAmount:       buyAmount,
		SellAmount:      sellAmount,
		Expires:         expires,
		Nonce:           big.NewInt(0),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		BuyToken:        buyTokenAddr,
		SellToken:       sellTokenAddr,
		UserAddress:     maker.Address,
	}

	order.Sign(maker)

	trade := &types.Trade{
		OrderHash:  order.Hash,
		Amount:     amount,
		Taker:      taker.Address,
		TradeNonce: big.NewInt(0),
	}

	trade.Sign(taker)

	exchange.SetTxSender(admin)
	_, err = exchange.Trade(order, trade)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	simulator.Commit()

	// TokenSell: InitialSellTokenAmount + amount * (amountSell/amountBuy)
	sellTokenTakerBalance, _ := sellToken.BalanceOf(taker.Address)
	sellTokenMakerBalance, _ := sellToken.BalanceOf(maker.Address)
	buyTokenTakerBalance, _ := buyToken.BalanceOf(taker.Address)
	buyTokenMakerBalance, _ := buyToken.BalanceOf(maker.Address)

	if sellTokenTakerBalance.Cmp(amount) != 0 {
		t.Errorf("Expected Taker balance of sellToken to be equal to %v but got %v instead", 5*1e17, sellTokenTakerBalance)
	}

	if sellTokenMakerBalance.Cmp(amount) != 0 {
		t.Errorf("Expected Maker balance of sellToken to be equal to %v but got %v instead", 5*1e17, sellTokenMakerBalance)
	}

	if buyTokenTakerBalance.Cmp(amount) != 0 {
		t.Errorf("Expected Taker balance of buyToken to be equal to %v but got %v instead", 5*1e17, buyTokenTakerBalance)
	}

	if buyTokenMakerBalance.Cmp(amount) != 0 {
		t.Errorf("Expected Taker balance of buyToken to be equal to %v but got %v instead", 5*1e17, buyTokenMakerBalance)
	}
}
