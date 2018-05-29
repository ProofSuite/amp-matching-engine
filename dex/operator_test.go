package dex

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

func TestSimulatedOperator(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	simulator := deployer.Backend.(*backends.SimulatedBackend)
	if err != nil {
		t.Errorf("Could not create default simulator: %v", err)
	}

	admin := config.Wallets[0]
	maker := config.Wallets[1]
	taker := config.Wallets[2]
	ZRXAmountInt := int64(1e18)
	WETHAmountInt := int64(1e18)
	ZRXAmount := big.NewInt(1e18)
	WETHAmount := big.NewInt(1e18)

	ZRXTokenContract, err := deployer.DeployToken(maker.Address, ZRXAmount)
	if err != nil {
		t.Errorf("Could not deploy ZRX: %v", err)
	}

	WETHTokenContract, err := deployer.DeployToken(taker.Address, WETHAmount)
	if err != nil {
		t.Errorf("Could not deploy WETH: %v", err)
	}

	ZRX := Token{Symbol: "ZRX", Address: ZRXTokenContract.Address}
	WETH := Token{Symbol: "WETH", Address: WETHTokenContract.Address}
	ZRXWETH := NewPair(ZRX, WETH)

	simulator.Commit()

	dex, err := deployer.DeployExchange(admin.Address)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	makerFactory := NewOrderFactory(&ZRXWETH, maker)
	takerFactory := NewOrderFactory(&ZRXWETH, taker)

	makerFactory.SetExchangeAddress(dex.Address)
	makerFactory.SetExchangeAddress(dex.Address)

	errChannel, err := dex.ListenToErrorEvents()
	if err != nil {
		t.Errorf("Could not get error channel: %v", err)
	}

	go func() {
		for {
			errLog := <-errChannel
			t.Errorf("New Error event: %v", errLog)
		}
	}()

	simulator.Commit()

	_, err = ZRXTokenContract.ApproveFrom(maker, dex.Address, ZRXAmount)
	if err != nil {
		t.Errorf("Could not approve token1: %v", err)
	}
	_, err = WETHTokenContract.ApproveFrom(taker, dex.Address, WETHAmount)
	if err != nil {
		t.Errorf("Could not approve token2: %v", err)
	}

	simulator.Commit()

	operator, err := NewOperatorFromContract(admin, dex, deployer.Backend)
	if err != nil {
		t.Errorf("Could not instantiate operator from contract: %v", err)
	}

	simulator.Commit()

	_, err = dex.DepositTokenFrom(maker, ZRX.Address, ZRXAmount)
	if err != nil {
		t.Errorf("Could not deposit token: %v", err)
	}

	_, err = dex.DepositTokenFrom(taker, WETH.Address, WETHAmount)
	if err != nil {
		t.Errorf("Could not deposit token: %v", err)
	}

	simulator.Commit()

	o, _ := makerFactory.NewOrder(WETH, WETHAmountInt, ZRX, ZRXAmountInt)
	trade, _ := takerFactory.NewTrade(o, WETHAmountInt)

	_, err = operator.ExecuteTrade(o, trade)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	simulator.Commit()

	TakerZRXBalance, _ := dex.TokenBalance(taker.Address, ZRX.Address)
	TakerWETHBalance, _ := dex.TokenBalance(taker.Address, WETH.Address)
	MakerZRXBalance, _ := dex.TokenBalance(maker.Address, ZRX.Address)
	MakerWETHBalance, _ := dex.TokenBalance(maker.Address, WETH.Address)

	if MakerZRXBalance.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", 0, MakerZRXBalance)
	}

	if TakerZRXBalance.Cmp(ZRXAmount) != 0 {
		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", ZRXAmount, TakerZRXBalance)
	}

	if TakerWETHBalance.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 0, TakerWETHBalance)
	}

	if MakerWETHBalance.Cmp(WETHAmount) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", WETHAmount, MakerWETHBalance)
	}

}
