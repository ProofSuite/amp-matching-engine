package contracts_test

import (
	"log"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
)

func SetupTokenTest() (*testutils.Deployer, *types.Wallet) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	_, err = daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	wallet := testutils.GetTestWallet()
	walletDao := new(mocks.WalletDao)
	walletDao.On("GetDefaultAdminWallet").Return(wallet, nil)

	walletService := services.NewWalletService(walletDao)
	txService := services.NewTxService(walletDao, wallet)

	deployer, err := testutils.NewSimulator(walletService, txService, []common.Address{wallet.Address})
	if err != nil {
		panic(err)
	}

	return deployer, wallet
}

func TestBalanceOf(t *testing.T) {
	deployer, _ := SetupTokenTest()

	receiver := testutils.GetTestAddress1()
	amount := big.NewInt(1e18)

	token, _, _, err := deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	balance, err := token.BalanceOf(receiver)
	if err != nil {
		t.Errorf("Error retrieving token balance: %v", err)
	}

	if balance.Cmp(amount) != 0 {
		t.Errorf("Token balance incorrect. Expected %v but instead got %v", amount, balance)
	}
}

func TestTotalSupply(t *testing.T) {
	deployer, _ := SetupTokenTest()

	receiver := testutils.GetTestAddress1()
	amount := big.NewInt(1e18)

	token, _, _, err := deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	supply, err := token.TotalSupply()
	if err != nil {
		t.Errorf("Error retrieving total supply")
	}

	if supply.Cmp(amount) != 0 {
		t.Errorf("Token Balance Incorrect. Expected %v but instead got %v", amount, supply)
	}
}

func TestTransfer(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	owner := wallet.Address
	receiver := testutils.GetTestAddress1()
	initialAmount := big.NewInt(1e18)
	transferAmount := big.NewInt(5e17)

	token, _, _, err := deployer.DeployToken(owner, initialAmount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	_, err = token.Transfer(receiver, transferAmount)
	if err != nil {
		t.Errorf("Could not transfer tokens: %v", err)
	}

	simulator.Commit()

	receiverBalance, err := token.BalanceOf(receiver)
	if err != nil {
		t.Errorf("Could not retrieve receiver balance %v", err)
	}
	if receiverBalance.Cmp(big.NewInt(5e17)) != 0 {
		t.Errorf("Expected receiver balance to be equal to 1/2e18 but got %v instead", receiverBalance)
	}
}

func TestApprove(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	owner := wallet.Address
	spender := testutils.GetTestAddress2()
	amount := big.NewInt(1e18)

	token, _, _, err := deployer.DeployToken(owner, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	_, err = token.Approve(spender, amount)
	if err != nil {
		t.Errorf("Could not approve tokens: %v", err)
	}

	simulator.Commit()

	allowance, err := token.Allowance(owner, spender)
	if err != nil {
		t.Errorf("Could not retrieve receiver allowance %v", err)
	}
	if allowance.Cmp(amount) != 0 {
		t.Errorf("Expected receiver balance to be equal to 1/2e18 but got %v instead", allowance)
	}
}

func TestTransferEvent(t *testing.T) {
	deployer, wallet := SetupTokenTest()

	owner := wallet.Address
	receiver := testutils.GetTestAddress2()

	logs := []*contractsinterfaces.TokenTransfer{}
	amount := big.NewInt(1e18)
	done := make(chan bool)

	token, _, _, err := deployer.DeployToken(owner, amount)
	if err != nil {
		t.Errorf("Could not deploy token: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	events, err := token.ListenToTransferEvents()
	if err != nil {
		t.Errorf("Could not open transfer events channel")
	}

	go func() {
		for {
			event := <-events
			logs = append(logs, event)
			done <- true
		}
	}()

	_, err = token.Transfer(receiver, amount)
	if err != nil {
		t.Errorf("Could not transfer tokens: %v", err)
	}

	simulator.Commit()
	<-done

	if len(logs) != 1 {
		t.Errorf("Events log has not the correct length")
	}

	parsedTransfer := logs[0]
	if parsedTransfer.From != owner {
		t.Errorf("Event 'From' field is not correct")
	}
	if parsedTransfer.To != receiver {
		t.Errorf("Event 'To' field is not correct")
	}
	if parsedTransfer.Value.Cmp(amount) != 0 {
		t.Errorf("Event 'Amount' field is not correct")
	}
}
