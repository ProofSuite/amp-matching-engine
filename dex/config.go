package dex

import (
	"fmt"
	"math/big"

	. "github.com/ethereum/go-ethereum/common"
)

var config = NewDefaultConfiguration()

// Config holds the general configuration for the application
// Currently some parameters are redundant as different type of
// configurations need different (but similar) parameters
// - Accounts is a list of ethereum address
// - Admin is the ethereum account that is in charge of sending transactions to the exchange contract
// - Wallets is a list of ethereum accounts that can be used for testing purposes. These accounts must be unlocked
// - Exchange is the address of the Exchange.sol smart contract
// - Tokens is a mapping of token symbols to token addresses (ERC-20 tokens)
// - QuoteTokens is a list of tokens that have been registered as quote currencies for testing purposes.
// - TokenPairs is a list of token pairs (base token + quote token) to be registered for testing purposes.
type Config struct {
	Accounts       []Address
	Admin          *Wallet
	Wallets        []*Wallet
	Constants      *Constants
	Exchange       Address
	Tokens         map[string]Address
	QuoteTokens    Tokens
	TokenPairs     TokenPairs
	Deployer       *Deployer
	OperatorParams *OperatorParams
}

type Constants struct {
	ether           *big.Int
	defaultGasPrice *big.Int
	defaultMaxGas   *big.Int
}

// NewDefaultConfiguration() returns the configuration used for testing. No contract
// is deployed and the data used does not exist on-chain.
func NewDefaultConfiguration() *Config {
	wallets := getWallets()
	accounts := getAccounts()
	tokens := getTokenAddresses()
	exchange := getExchangeAddress()
	quoteTokens := getQuoteTokens()
	tokenPairs := getTokenPairs()

	admin := wallets[0]

	return &Config{
		Admin:       admin,
		Exchange:    exchange,
		Wallets:     wallets,
		Accounts:    accounts,
		Tokens:      tokens,
		QuoteTokens: quoteTokens,
		TokenPairs:  tokenPairs,
	}
}

// NewSimulatorConfiguration() returns the configuration used for testing. Contracts
// are deployed on a go-ethereum simulated backend.
// ZRX, EOS, WETH tokens (which are actually only standard ERC20 tokens arbitrarily named
// for clarity purposes) are deployed and then added to the pair, quotes and token variables.
func NewSimulatorConfiguration() *Config {
	minted := big.NewInt(1e18)
	constants := getConstants()
	wallets := getWallets()
	accounts := getAccounts()

	admin := wallets[0]

	deployer, err := NewSimulator(admin, accounts)
	if err != nil {
		fmt.Printf("Could not deploy simulator")
	}

	ZRXTokenContract, _, err := deployer.DeployToken(admin.Address, minted)
	if err != nil {
		fmt.Printf("Could not deploy the ZRX token contract")
	}

	WETHTokenContract, _, err := deployer.DeployToken(admin.Address, minted)
	if err != nil {
		fmt.Printf("Could not deploy the WETH token contract")
	}

	EOSTokenContract, _, err := deployer.DeployToken(admin.Address, minted)
	if err != nil {
		fmt.Printf("Could not deploy the EOS token contract")
	}

	ZRX := Token{Symbol: "ZRX", Address: ZRXTokenContract.Address}
	WETH := Token{Symbol: "WETH", Address: WETHTokenContract.Address}
	EOS := Token{Symbol: "EOS", Address: EOSTokenContract.Address}
	ZRX_WETH := NewPair(ZRX, WETH)
	EOS_WETH := NewPair(EOS, WETH)

	pairs := TokenPairs{}
	pairs["ZRXWETH"] = ZRX_WETH
	pairs["EOSWETH"] = EOS_WETH

	quotes := Tokens{}
	quotes["ZRX"] = ZRX
	quotes["WETH"] = WETH
	quotes["EOS"] = EOS

	tokens := make(map[string]Address)
	tokens["ZRX"] = ZRX.Address
	tokens["WETH"] = WETH.Address
	tokens["EOS"] = EOS.Address

	ex, _, err := deployer.DeployExchange(admin.Address)
	if err != nil {
		fmt.Printf("Could not deploy exchange: %v", err)
	}

	return &Config{
		Admin:       admin,
		Constants:   constants,
		Exchange:    ex.Address,
		Wallets:     wallets,
		Accounts:    accounts,
		Tokens:      tokens,
		QuoteTokens: quotes,
		TokenPairs:  pairs,
		Deployer:    deployer,
	}
}

// CreateLocalhostConfiguration() deploys mock tokens and the decentralized exchange contract
// and returns a configuration object that includes accounts, contract addresses, etc.
// If contracts are already deployed, then use the NewLocalHostConfiguration instead.
func CreateConfiguration() *Config {
	minted := big.NewInt(1e18)
	constants := getConstants()
	wallets := getWallets()
	accounts := getAccounts()
	opParams := getOperatorParams()
	admin := wallets[0]

	deployer, err := NewDeployer(admin)
	if err != nil {
		fmt.Printf("Could not deploy simulator")
	}

	ex, _, err := deployer.DeployExchange(admin.Address)
	if err != nil {
		fmt.Printf("Could not deploy the exchange contract")
	}

	ZRXContract, _, err := deployer.DeployToken(admin.Address, minted)
	if err != nil {
		fmt.Printf("Could not deploy the ZRX token contract")
	}

	WETHContract, _, err := deployer.DeployToken(admin.Address, minted)
	if err != nil {
		fmt.Printf("Could not deploy the WETH token contract")
	}

	EOSContract, _, err := deployer.DeployToken(admin.Address, minted)
	if err != nil {
		fmt.Printf("Could not deploy the EOS token contract")
	}

	ZRX := Token{Symbol: "ZRX", Address: ZRXContract.Address}
	WETH := Token{Symbol: "WETH", Address: WETHContract.Address}
	EOS := Token{Symbol: "EOS", Address: EOSContract.Address}
	ZRX_WETH := NewPair(ZRX, WETH)
	EOS_WETH := NewPair(EOS, WETH)

	pairs := TokenPairs{}
	pairs["ZRXWETH"] = ZRX_WETH
	pairs["EOSWETH"] = EOS_WETH

	quotes := Tokens{}
	quotes["ZRX"] = ZRX
	quotes["WETH"] = WETH
	quotes["EOS"] = EOS

	tokens := make(map[string]Address)
	tokens["ZRX"] = ZRX.Address
	tokens["WETH"] = WETH.Address
	tokens["EOS"] = EOS.Address

	return &Config{
		Admin:          admin,
		Constants:      constants,
		Exchange:       ex.Address,
		Wallets:        wallets,
		Accounts:       accounts,
		Tokens:         tokens,
		QuoteTokens:    quotes,
		TokenPairs:     pairs,
		Deployer:       deployer,
		OperatorParams: opParams,
	}

}

// NewLocalhostConfiguration() returns the configuration used for interacting with a localhost
// blockchain. ZRX, EOS, WETH token objects are created from the given addresses.
// The DEX, ZRX, EOS, WETH contracts (ZRX, EOS, WETH can be any ERC20 contracts) have to be
// already deployed. Otherwise, use the CreateLocalHostConfiguration contract.
func NewConfiguration() *Config {
	constants := getConstants()
	wallets := getWallets()
	accounts := getAccounts()
	tokens := getTokenAddresses()
	exchange := getExchangeAddress()
	opParams := getOperatorParams()

	admin := wallets[0]

	dep, err := NewDeployer(admin)
	if err != nil {
		fmt.Printf("Could not deploy simulator")
	}

	ZRXContract, err := dep.NewToken(tokens["ZRX"])
	if err != nil {
		fmt.Printf("Could not retrieve the ZRX token contract")
	}

	WETHContract, err := dep.NewToken(tokens["WETH"])
	if err != nil {
		fmt.Printf("Could not retrieve the WETH token contract")
	}

	EOSContract, err := dep.NewToken(tokens["EOS"])
	if err != nil {
		fmt.Printf("Could not retrieve the EOS token contract")
	}

	ex, err := dep.NewExchange(exchange)
	if err != nil {
		fmt.Printf("Could not deploy the exchange contract")
	}

	ZRX := Token{Symbol: "ZRX", Address: ZRXContract.Address}
	WETH := Token{Symbol: "WETH", Address: WETHContract.Address}
	EOS := Token{Symbol: "EOS", Address: EOSContract.Address}
	ZRX_WETH := NewPair(ZRX, WETH)
	EOS_WETH := NewPair(EOS, WETH)

	pairs := TokenPairs{}
	pairs["ZRXWETH"] = ZRX_WETH
	pairs["EOSWETH"] = EOS_WETH

	quotes := Tokens{}
	quotes["ZRX"] = ZRX
	quotes["WETH"] = WETH
	quotes["EOS"] = EOS

	return &Config{
		Admin:          admin,
		Constants:      constants,
		Exchange:       ex.Address,
		Wallets:        wallets,
		Accounts:       accounts,
		Tokens:         tokens,
		QuoteTokens:    quotes,
		TokenPairs:     pairs,
		Deployer:       dep,
		OperatorParams: opParams,
	}
}

// getTokens returns a mapping of token symbol to token addresses
// These addresses correspond to the addresses of tokens that have been deployed on the private ethereum
// network which you can find a repository on github.com/amansardana/private-geth-chain
func getTokenAddresses() map[string]Address {
	tokens := make(map[string]Address)
	tokens["EOS"] = HexToAddress("0x5d564669ab4cfd96b785d3d05e8c7d66a073daf0")
	tokens["ZRX"] = HexToAddress("0x9792845456a0075df8a03123e7dac62bb0f69440")
	tokens["WETH"] = HexToAddress("0x27cb1d4b335ec45512088eea990238344d776714")

	return tokens
}

// getExchangeAddress returns the address that have been deployed on the private ethereum network that
// you can find on github.com/amansardana/private-geth-chain
func getExchangeAddress() Address {
	return HexToAddress("0x29faee20f205c15c6c3004482f8996a468336b67")
}

// getAccounts returns default accounts. These addresses can be funded and unlocked on the private ethereum
// network that you can find on github.com/amansardana/private-geth-chain
func getAccounts() []Address {
	accountList := []Address{}

	addresses := []string{
		"0xe8e84ee367bc63ddb38d3d01bccef106c194dc47",
		"0xcf7389dc6c63637598402907d5431160ec8972a5",
		"0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa",
		"0x14d281013d8ee8ccfa0eca87524e5b3cfa6152ba",
		"0x6e9a406696617ec5105f9382d33ba3360fcfabcc",
		"0x7e0f08462bf391ee4154a88994f8ce2aad7ab144",
		"0x4dc5790733b997f3db7fc49118ab013182d6ba9b",
		"0x545aeb22f378ef7a4f627c45efe8245152bed8a1",
		"0x830212529506afd9c24adcfdde6fe825982d37ae",
		"0x44809695706c252435531029b1e9d7d0355d475f",
	}

	for _, address := range addresses {
		accountList = append(accountList, HexToAddress(address))
	}

	return accountList
}

// getConstants returns ether value in wei, the default gas price in wei
// and the default max gas in wei
func getConstants() *Constants {
	ether := big.NewInt(1e18)
	defaultGasPrice := big.NewInt(1e9)
	defaultMaxGas := big.NewInt(5e6)

	return &Constants{
		ether:           ether,
		defaultGasPrice: defaultGasPrice,
		defaultMaxGas:   defaultMaxGas,
	}
}

// getWallets returns a list of private keys that can be funded on the private ehtereum network
// that you can find on github.com/amansardana/private-geth-chain
func getWallets() []*Wallet {
	walletList := []*Wallet{}

	keys := []string{
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712661",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712662",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712663",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712664",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712665",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712666",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712667",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712668",
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712669",
	}

	for _, key := range keys {
		walletList = append(walletList, NewWalletFromPrivateKey(key))
	}

	return walletList
}

// getQuoteTokens generates a mapping of quote tokens (arbitrary addresses). These addresses correspond to
// tokens that have been deployed on the github.com/amansardana/private-geth-chain private ethereum network
func getQuoteTokens() Tokens {
	quoteTokens := Tokens{}

	WETH := Token{Symbol: "WETH", Address: HexToAddress("0x5d564669ab4cfd96b785d3d05e8c7d66a073daf0")}
	ZRX := Token{Symbol: "ZRX", Address: HexToAddress("0x9792845456a0075df8a03123e7dac62bb0f69440")}
	EOS := Token{Symbol: "EOS", Address: HexToAddress("0x27cb1d4b335ec45512088eea990238344d776714")}

	quoteTokens["WETH"] = WETH
	quoteTokens["ZRX"] = ZRX
	quoteTokens["EOS"] = EOS
	return quoteTokens
}

// getTokenPairs generates a mapping of token pairs (arbitrary addresses). These addresses correspond to
// tokens that have been deployed on the github.com/amansardana/private-geth-chain private ethereum network
func getTokenPairs() TokenPairs {
	tokenPairs := TokenPairs{}

	WETH := Token{Symbol: "WETH", Address: HexToAddress("0x5d564669ab4cfd96b785d3d05e8c7d66a073daf0")}
	ZRX := Token{Symbol: "ZRX", Address: HexToAddress("0x9792845456a0075df8a03123e7dac62bb0f69440")}
	EOS := Token{Symbol: "EOS", Address: HexToAddress("0x27cb1d4b335ec45512088eea990238344d776714")}

	ZRX_WETH := TokenPair{BaseToken: ZRX, QuoteToken: WETH}
	EOS_WETH := TokenPair{BaseToken: EOS, QuoteToken: WETH}

	ZRX_WETH = NewPair(ZRX, WETH)
	EOS_WETH = NewPair(EOS, WETH)

	tokenPairs["ZRXWETH"] = ZRX_WETH
	tokenPairs["EOSWETH"] = EOS_WETH
	return tokenPairs
}

func getOperatorParams() *OperatorParams {
	return &OperatorParams{
		gasPrice:   big.NewInt(1e9),
		maxGas:     5e6,
		minBalance: big.NewInt(1e17),
		rpcURL:     "ws://127.0.0.1:8546",
	}
}
