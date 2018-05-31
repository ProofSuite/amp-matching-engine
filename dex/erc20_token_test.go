package dex

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Dvisacker/matching-engine/dex/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

func TestBalanceOf(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		fmt.Printf("%v", err)
	}
	receiver := config.Accounts[1]
	amount := big.NewInt(1e18)

	token, _, err := deployer.DeployToken(receiver, amount)
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
	deployer, _ := NewDefaultSimulator()
	receiver := config.Accounts[0]
	amount := big.NewInt(1e18)

	token, _, err := deployer.DeployToken(receiver, amount)
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
	deployer, _ := NewDefaultSimulator()
	owner := config.Accounts[0]
	receiver := config.Accounts[1]
	initialAmount := big.NewInt(1e18)
	transferAmount := big.NewInt(1e18 / 2)

	token, _, err := deployer.DeployToken(owner, initialAmount)
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
	if receiverBalance.Cmp(big.NewInt(1e18/2)) != 0 {
		t.Errorf("Expected receiver balance to be equal to 1/2e18 but got %v instead", receiverBalance)
	}
}

func TestApprove(t *testing.T) {
	deployer, _ := NewDefaultSimulator()
	owner := config.Accounts[0]
	spender := config.Accounts[1]
	amount := big.NewInt(1e18)

	token, _, err := deployer.DeployToken(owner, amount)
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
	deployer, _ := NewDefaultSimulator()
	logs := []*interfaces.TokenTransfer{}
	owner := config.Accounts[0]
	receiver := config.Accounts[1]
	amount := big.NewInt(1e18)
	done := make(chan bool)

	token, _, err := deployer.DeployToken(owner, amount)
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
