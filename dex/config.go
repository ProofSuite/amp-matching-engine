package dex

import (
	. "github.com/ethereum/go-ethereum/common"
)

var config = NewDefaultConfiguration()

type Config struct {
	Contracts   ConfigContracts
	Wallets     []*Wallet
	Keys        []string
	Accounts    []Address
	Tokens      map[string]Address
	QuoteTokens Tokens
	TokenPairs  TokenPairs
}

type ConfigContracts struct {
	exchange Address
	token1   Address
	token2   Address
}

func NewDefaultConfiguration() *Config {
	wallets := getWallets()
	accounts := getAccounts()
	tokens := getTokens()
	contracts := getContracts()
	quoteTokens := getQuoteTokens()
	tokenPairs := getTokenPairs()

	return &Config{
		Contracts: contracts,
		Wallets:   wallets,
		Accounts:  accounts,

		Keys: []string{
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
		},
		Tokens:      tokens,
		QuoteTokens: quoteTokens,
		TokenPairs:  tokenPairs,
	}
}

func getTokens() map[string]Address {
	tokens := make(map[string]Address)
	tokens["EOS"] = HexToAddress("0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0")
	tokens["ZRX"] = HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")

	return tokens
}

func getContracts() ConfigContracts {
	return ConfigContracts{
		exchange: HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		token1:   HexToAddress("0x5ac05570112c0a95f2fd1d85292e0c08522a1bdb"),
		token2:   HexToAddress("0x43925198636a2d43f1c887ed1d936c76f03f55de"),
	}
}

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

// func getQuoteTokens() []*Token {

// 	EOS := &Token{Symbol: "EOS", Address: HexToAddress("0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0")}
// 	WETH := &Token{Symbol: "WETH", Address: HexToAddress("0x2956356cd2a2bf3202f771f50d3d14a367b48070")}

// 	quoteTokens := []*Token{EOS, WETH}
// 	return quoteTokens
// }

func getQuoteTokens() Tokens {
	quoteTokens := Tokens{}

	WETH := Token{Symbol: "WETH", Address: HexToAddress("0x2956356cd2a2bf3202f771f50d3d14a367b48070")}
	ZRX := Token{Symbol: "ZRX", Address: HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")}
	EOS := Token{Symbol: "EOS", Address: HexToAddress("0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0")}

	quoteTokens["WETH"] = WETH
	quoteTokens["ZRX"] = ZRX
	quoteTokens["EOS"] = EOS
	return quoteTokens
}

func getTokenPairs() TokenPairs {
	tokenPairs := TokenPairs{}

	WETH := Token{Symbol: "WETH", Address: HexToAddress("0x2956356cd2a2bf3202f771f50d3d14a367b48070")}
	ZRX := Token{Symbol: "ZRX", Address: HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")}
	EOS := Token{Symbol: "EOS", Address: HexToAddress("0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0")}

	ZRX_WETH := TokenPair{BaseToken: ZRX, QuoteToken: WETH}
	EOS_WETH := TokenPair{BaseToken: EOS, QuoteToken: WETH}

	ZRX_WETH = NewPair(ZRX, WETH)
	EOS_WETH = NewPair(EOS, WETH)

	tokenPairs["ZRXWETH"] = ZRX_WETH
	tokenPairs["EOSWETH"] = EOS_WETH
	return tokenPairs
}
