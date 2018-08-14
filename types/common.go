package types

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewAddressFromString(addr string) common.Address {
	return common.HexToAddress(addr)
}

func NewPrivateKeyFromString(key string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func NewBigInt(n string) (bn *big.Int) {
	bn = new(big.Int)
	bn, _ = bn.SetString(n, 10)
	return
}
