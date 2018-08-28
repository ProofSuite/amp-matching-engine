package testutils

import "github.com/Proofsuite/amp-matching-engine/types"

func GetTestWallet() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
}

func GetTestWallet1() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
}

func GetTestWallet2() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712661")
}

func GetTestWallet3() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712662")
}
