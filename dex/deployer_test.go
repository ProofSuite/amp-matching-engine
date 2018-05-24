package dex

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
)

func TestNewDefaultLocalDeployer(t *testing.T) {
	_, err := NewDefaultLocalDeployer()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}
}

func TestNewDefaultSimulator(t *testing.T) {
	_, err := NewDefaultSimulator()
	if err != nil {
		t.Errorf("Error creating simulator: %v", err)
	}
}

func TestDeployTokenWithLocalBackend(t *testing.T) {
	deployer, err := NewDefaultLocalDeployer()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	wallet := deployer.Wallet
	receiver := wallet.Address
	amount := big.NewInt(1e18)

	_, err = deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Error deploying token: %v", err)
	}
}

func TestNewTokenWithLocalBackend(t *testing.T) {
	deployer, err := NewDefaultLocalDeployer()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	address := config.Contracts.token1
	_, err = deployer.NewToken(address)
	if err != nil {
		t.Errorf("Error deploying token: %v", err)
	}
}

func TestDeployTokenWithSimulatedBackend(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	wallet := deployer.Wallet
	receiver := wallet.Address
	amount := big.NewInt(1e18)

	_, err = deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Error deploying token: %v", err)
	}
}

func TestNewTokenWithSimulatedBackend(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	wallet := deployer.Wallet
	receiver := wallet.Address
	amount := big.NewInt(1e18)

	_, err = deployer.DeployToken(receiver, amount)
	if err != nil {
		t.Errorf("Error deploying token")
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

}

func TestDeployExchangeWithLocalBackend(t *testing.T) {
	deployer, err := NewDefaultLocalDeployer()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	feeAccount := config.Accounts[1]

	_, err = deployer.DeployExchange(feeAccount)
	if err != nil {
		t.Errorf("Error deploying token: %v", err)
	}
}

func TestNewExchangeWithLocalBackend(t *testing.T) {
	deployer, err := NewDefaultLocalDeployer()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	address := config.Contracts.exchange
	_, err = deployer.NewExchange(address)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}
}

func TestDeployExchangeWithSimulatedBackend(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	feeAccount := config.Accounts[1]

	_, err = deployer.DeployExchange(feeAccount)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}
}

func TestNewExchangeWithSimulatedBackend(t *testing.T) {
	deployer, err := NewDefaultSimulator()
	if err != nil {
		t.Errorf("Error creating deployer: %v", err)
	}

	feeAccount := config.Accounts[1]

	exchange, err := deployer.NewExchange(feeAccount)
	if err != nil {
		t.Errorf("Error deploying exchange: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()
	// deployer.(*backend).Commit()

	_, err = deployer.NewExchange(exchange.Address)
	if err != nil {
		t.Errorf("Error getting new exchange contract instance: %v", err)
	}
}
