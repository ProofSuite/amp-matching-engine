package dex

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

func TestSetFeeAccount(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	initialFeeAccount := config.Wallets[0].Address
	newFeeAccount := config.Wallets[1].Address

	exchange, err := deployer.DeployExchange(initialFeeAccount)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	_, err = exchange.SetFeeAccount(newFeeAccount)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	simulator.Commit()

	feeAccount, err := exchange.FeeAccount()
	if err != nil {
		t.Errorf("Error retrieving fee account address: %v", err)
	}

	if newFeeAccount != feeAccount {
		t.Errorf("Fee account not set correctly")
	}
}

func TestSetOperator(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	feeAccount := config.Wallets[0].Address
	operator := config.Wallets[1].Address
	exchange, err := deployer.DeployExchange(feeAccount)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

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

func TestSetWithdrawalSecurityPeriod(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	feeAccount := config.Wallets[0].Address
	exchange, err := deployer.DeployExchange(feeAccount)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	period := big.NewInt(100000000)
	_, err = exchange.SetWithdrawalSecurityPeriod(period)
	if err != nil {
		t.Errorf("Error calling the operator variable: %v", err)
	}

	newWithdrawalPeriod, err := exchange.WithdrawalSecurityPeriod()
	if err != nil {
		t.Errorf("Error getting the withdrawal period")
	}

	if newWithdrawalPeriod.Cmp(period) != 0 {
		t.Errorf("Expected withdrawal period to be equal to %v but instead got %v", period, newWithdrawalPeriod)
	}
}

func TestDepositEther(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	feeAccount := config.Wallets[0].Address
	exchange, err := deployer.DeployExchange(feeAccount)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	sender := config.Wallets[1]
	value := big.NewInt(5 * 1e17)

	exchange.SetCustomSender(sender)
	_, err = exchange.DepositEther(value)
	if err != nil {
		t.Errorf("Could not deposit ether: %v", err)
	}

	simulator.Commit()

	balance, err := exchange.EtherBalance(sender.Address)
	if err != nil {
		t.Errorf("Error retrieving ether balance: %v", err)
	}
	if balance.Cmp(value) != 0 {
		t.Errorf("Balance error: Expected %v but instead got %v", value, balance)
	}
}

func TestDepositToken(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	admin := config.Wallets[0].Address
	sender := config.Wallets[1]
	simulator := deployer.Backend.(*backends.SimulatedBackend)
	amount := big.NewInt(1e18)

	exchange, err := deployer.DeployExchange(admin)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}

	token, err := deployer.DeployToken(sender.Address, amount)
	if err != nil {
		t.Errorf("Error deploying token: %v", err)
	}

	simulator.Commit()

	token.SetCustomSender(sender)
	_, err = token.Approve(exchange.Address, amount)
	if err != nil {
		t.Errorf("Could not approve token: %v", err)
	}

	simulator.Commit()

	exchange.SetCustomSender(sender)
	_, err = exchange.DepositToken(token.Address, amount)
	if err != nil {
		t.Errorf("Could not deposit token: %v", err)
	}

	simulator.Commit()

	balance, err := exchange.TokenBalance(sender.Address, token.Address)
	if err != nil {
		t.Errorf("Could not retrieve token balance: %v", err)
	}
	if balance.Cmp(amount) != 0 {
		t.Errorf("Balance error: Expected %v but instead got %v", amount, balance)
	}
}

func TestWithdraw(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	admin := config.Wallets[0]
	sender := config.Wallets[1]
	receiver := config.Wallets[2]
	amount := big.NewInt(1e18)
	n := big.NewInt(1)
	f := big.NewInt(0)

	exchange, err := deployer.DeployExchange(admin.Address)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}

	token, err := deployer.DeployToken(sender.Address, amount)
	if err != nil {
		t.Errorf("Error deploying token: %v", err)
	}

	simulator.Commit()

	token.SetCustomSender(sender)
	_, err = token.Approve(exchange.Address, amount)
	if err != nil {
		t.Errorf("Could not approve token: %v", err)
	}

	simulator.Commit()

	exchange.SetCustomSender(sender)
	_, err = exchange.DepositToken(token.Address, amount)
	if err != nil {
		t.Errorf("Could not deposit token: %v", err)
	}

	simulator.Commit()

	w := &Withdrawal{
		ExchangeAddress: exchange.Address,
		Token:           token.Address,
		Amount:          amount,
		Trader:          sender.Address,
		Receiver:        receiver.Address,
		Nonce:           n,
		Fee:             f,
	}

	err = w.Sign(sender)
	if err != nil {
		t.Errorf("Could not sign withdrawal payload: %v", err)
	}

	//The exchange operator is the only address allowed to perform withdraws
	exchange.SetDefaultSender()
	_, err = exchange.Withdraw(w)
	if err != nil {
		t.Errorf("Could not deposit token: %v", err)
	}

	simulator.Commit()

	//The receiver address is the address to which the tokens are withdrawn
	balance, err := token.BalanceOf(receiver.Address)
	if err != nil {
		t.Errorf("Could not retrieve receiver token balance: %v", err)
	}

	//The sender/trader address is the address at which the tokens are originally
	//attributed on the exchange.
	exchangeTokenBalance, err := exchange.TokenBalance(token.Address, sender.Address)
	if err != nil {
		t.Errorf("Could not retrieve trader token balance")
	}

	if balance.Cmp(amount) != 0 {
		t.Errorf("Expected balance to be equal to %v but got %v", amount, balance)
	}

	if exchangeTokenBalance.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected balance to be equal to %v but got %v", amount, balance)
	}
}

func TestTrade(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)

	admin := config.Wallets[0]
	maker := config.Wallets[1]
	taker := config.Wallets[2]
	buyAmount := big.NewInt(1e18)
	sellAmount := big.NewInt(1e18)
	amount := big.NewInt(5e17)
	expires := big.NewInt(1e7)

	exchange, err := deployer.DeployExchange(admin.Address)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}

	//Initially Maker owns 1e18 units of sellToken and Taker owns 1e18 units buyToken
	sellToken, err := deployer.DeployToken(maker.Address, sellAmount)
	if err != nil {
		t.Errorf("Error deploying token 1: %v", err)
	}
	buyToken, err := deployer.DeployToken(taker.Address, buyAmount)
	if err != nil {
		t.Errorf("Error deploying token 2: %v", err)
	}

	simulator.Commit()
	exchange.PrintErrors()

	sellToken.SetCustomSender(maker)
	_, err = sellToken.Approve(exchange.Address, sellAmount)
	if err != nil {
		t.Errorf("Could not approve sellToken: %v", err)
	}

	buyToken.SetCustomSender(taker)
	_, err = buyToken.Approve(exchange.Address, buyAmount)
	if err != nil {
		t.Errorf("Could not approve buyToken: %v", err)
	}

	simulator.Commit()

	exchange.SetCustomSender(taker)
	_, err = exchange.DepositToken(buyToken.Address, buyAmount)
	if err != nil {
		t.Errorf("Could not deposit buyToken: %v", err)
	}

	exchange.SetCustomSender(maker)
	_, err = exchange.DepositToken(sellToken.Address, sellAmount)
	if err != nil {
		t.Errorf("Could not deposit token: %v", err)
	}

	simulator.Commit()

	//Maker creates an order that exchanges 'sellAmount' of sellToken for 'buyAmount' of buyToken
	order := &Order{
		ExchangeAddress: exchange.Address,
		AmountBuy:       buyAmount,
		AmountSell:      sellAmount,
		Expires:         expires,
		Nonce:           big.NewInt(0),
		FeeMake:         big.NewInt(0),
		FeeTake:         big.NewInt(0),
		TokenBuy:        buyToken.Address,
		TokenSell:       sellToken.Address,
		Maker:           maker.Address,
	}
	order.Sign(maker)

	trade := &Trade{
		OrderHash:  order.Hash,
		Amount:     amount,
		Taker:      taker.Address,
		TradeNonce: big.NewInt(0),
	}
	trade.Sign(taker)

	exchange.SetDefaultSender()
	_, err = exchange.Trade(order, trade)
	if err != nil {
		t.Errorf("Could not do execute trade: %v", err)
	}

	simulator.Commit()

	// At the end of the trade, without the fees into account
	// Maker should own:
	// TokenBuy: InitialBuyTokenAmount + amount
	// TokenSell: InitialSellTokenAmount - amount * (amountSell/amountBuy)

	// Taker should own:
	// TokenBuy: InitialBuyTokenAmount - amount
	// TokenSell: InitialSellTokenAmount + amount * (amountSell/amountBuy)
	sellTokenTakerBalance, _ := exchange.TokenBalance(taker.Address, sellToken.Address)
	sellTokenMakerBalance, _ := exchange.TokenBalance(maker.Address, sellToken.Address)
	buyTokenTakerBalance, _ := exchange.TokenBalance(taker.Address, buyToken.Address)
	buyTokenMakerBalance, _ := exchange.TokenBalance(maker.Address, buyToken.Address)

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
