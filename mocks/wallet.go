package mocks

import "github.com/Proofsuite/amp-matching-engine/types"

func getMockWallet() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
}
