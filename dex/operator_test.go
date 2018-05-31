package dex

import (
	"math/big"
	"testing"
)

func TestOperator(t *testing.T) {
	testConfig := NewConfiguration()
	opParams := testConfig.OperatorParams

	ZRX := testConfig.QuoteTokens["ZRX"]
	WETH := testConfig.QuoteTokens["WETH"]
	ZRXWETH := NewPair(ZRX, WETH)

	admin := config.Wallets[0]
	maker := config.Wallets[1]
	taker := config.Wallets[2]
	exchange := config.Exchange

	opConfig := &OperatorConfig{
		Admin:          admin,
		Exchange:       exchange,
		OperatorParams: opParams,
	}

	ZRXAmount := big.NewInt(1e18)
	WETHAmount := big.NewInt(1e18)

	deployer, err := NewWebsocketDeployer(admin)
	if err != nil {
		t.Errorf("Could not instantiate deployer: %v", err)
	}

	ex, err := deployer.NewExchange(exchange)
	if err != nil {
		t.Errorf("Could not retrieve exchange instance: %v", err)
	}

	makerFactory := NewOrderFactory(&ZRXWETH, maker)
	takerFactory := NewOrderFactory(&ZRXWETH, taker)
	makerFactory.SetExchangeAddress(ex.Address)
	takerFactory.SetExchangeAddress(ex.Address)

	operator, err := NewOperator(opConfig)
	if err != nil {
		t.Errorf("Could not instantiate operator: %v", err)
	}

	initialTakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
	initialTakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
	initialMakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
	initialMakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

	o, _ := makerFactory.NewOrder(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
	trade, _ := takerFactory.NewTrade(o, 1)

	tx, err := operator.ExecuteTrade(o, trade)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	_, err = operator.WaitMined(tx)
	if err != nil {
		t.Errorf("Could not mine trade transaction")
	}

	TakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
	TakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
	MakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
	MakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

	TakerZRXIncrement := big.NewInt(0)
	TakerWETHIncrement := big.NewInt(0)
	MakerZRXIncrement := big.NewInt(0)
	MakerWETHIncrement := big.NewInt(0)

	MakerZRXIncrement.Sub(MakerZRXBalance, initialMakerZRXBalance)
	MakerWETHIncrement.Sub(MakerWETHBalance, initialMakerWETHBalance)
	TakerZRXIncrement.Sub(TakerZRXBalance, initialTakerZRXBalance)
	TakerWETHIncrement.Sub(TakerWETHBalance, initialTakerWETHBalance)

	if MakerZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, MakerZRXIncrement)
	}

	if MakerWETHIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, MakerWETHIncrement)
	}

	if TakerWETHIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, TakerWETHIncrement)
	}

	if TakerZRXIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, TakerZRXIncrement)
	}
}

// func TestOperator2(t *testing.T) {

// 	testConfig := NewConfiguration()
// 	opParams := testConfig.OperatorParams

// 	admin := config.Wallets[0]
// 	maker := config.Wallets[1]
// 	taker := config.Wallets[2]
// 	exchange := config.Exchange

// 	opConfig := &OperatorConfig{
// 		Admin:          admin,
// 		Exchange:       exchange,
// 		OperatorParams: opParams,
// 	}

// 	ZRXAmount := big.NewInt(1e18)
// 	WETHAmount := big.NewInt(1e18)

// 	deployer, err := NewWebsocketDeployer(admin)
// 	if err != nil {
// 		t.Errorf("Could not instantiate deployer: %v", err)
// 	}
// 	ZRXTokenContract, tx, err := deployer.DeployToken(maker.Address, ZRXAmount)
// 	if err != nil {
// 		t.Errorf("Could not deploy ZRX: %v", err)
// 	}
// 	_, err = deployer.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine ZRX deployment transaction")
// 	}

// 	WETHTokenContract, tx, err := deployer.DeployToken(taker.Address, WETHAmount)
// 	if err != nil {
// 		t.Errorf("Could not deploy WETH: %v", err)
// 	}
// 	_, err = deployer.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine WETH deployment transaction")
// 	}

// 	ZRX := Token{Symbol: "ZRX", Address: ZRXTokenContract.Address}
// 	WETH := Token{Symbol: "WETH", Address: WETHTokenContract.Address}
// 	ZRXWETH := NewPair(ZRX, WETH)

// 	ex, tx, err := deployer.DeployExchange(admin.Address)
// 	if err != nil {
// 		t.Errorf("Could not deploy exchange: %v", err)
// 	}
// 	_, err = deployer.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine Exchange deployment transaction")
// 	}

// 	makerFactory := NewOrderFactory(&ZRXWETH, maker)
// 	takerFactory := NewOrderFactory(&ZRXWETH, taker)

// 	makerFactory.SetExchangeAddress(ex.Address)
// 	takerFactory.SetExchangeAddress(ex.Address)

// 	tx, err = ZRXTokenContract.ApproveFrom(maker, ex.Address, ZRXAmount)
// 	if err != nil {
// 		t.Errorf("Could not approve ZRX Token: %v", err)
// 	}
// 	_, err = deployer.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine approx ZRX Token Transaction")
// 	}

// 	tx, err = WETHTokenContract.ApproveFrom(taker, ex.Address, WETHAmount)
// 	if err != nil {
// 		t.Errorf("Could not approve WETH Token: %v", err)
// 	}
// 	_, err = deployer.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine approve WETH Token Transaction")
// 	}

// 	operator, err := NewOperator(opConfig)
// 	if err != nil {
// 		t.Errorf("Could not instantiate operator: %v", err)
// 	}

// 	tx, err = ex.DepositTokenFrom(maker, ZRX.Address, ZRXAmount)
// 	if err != nil {
// 		t.Errorf("Could not deposit token: %v", err)
// 	}
// 	_, err = operator.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine ZRX Token Deposit transaction")
// 	}

// 	tx, err = ex.DepositTokenFrom(taker, WETH.Address, WETHAmount)
// 	if err != nil {
// 		t.Errorf("Could not deposit token: %v", err)
// 	}
// 	_, err = operator.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine WETH Token Deposit transaction")
// 	}

// 	o, _ := makerFactory.NewOrder(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
// 	trade, _ := takerFactory.NewTrade(o, WETHAmount.Int64())

// 	tx, err = operator.ExecuteTrade(o, trade)
// 	if err != nil {
// 		t.Errorf("Could not execute trade: %v", err)
// 	}
// 	_, err = operator.WaitMined(tx)
// 	if err != nil {
// 		t.Errorf("Could not mine trade transaction")
// 	}

// 	TakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
// 	TakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
// 	MakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
// 	MakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

// 	if MakerZRXBalance.Cmp(big.NewInt(0)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", 0, MakerZRXBalance)
// 	}

// 	if TakerZRXBalance.Cmp(ZRXAmount) != 0 {
// 		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", ZRXAmount, TakerZRXBalance)
// 	}

// 	if TakerWETHBalance.Cmp(big.NewInt(0)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 0, TakerWETHBalance)
// 	}

// 	if MakerWETHBalance.Cmp(WETHAmount) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", WETHAmount, MakerWETHBalance)
// 	}
// }
