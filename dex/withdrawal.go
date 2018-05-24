package dex

import (
	"math/big"

	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Withdrawal struct {
	ExchangeAddress Address
	Hash            Hash
	Token           Address
	Amount          *big.Int
	Trader          Address
	Receiver        Address
	Nonce           *big.Int
	Fee             *big.Int
	Signature       *Signature
}

func (w *Withdrawal) ComputeWithdrawalHash() Hash {
	sha := sha3.NewKeccak256()

	sha.Write(w.ExchangeAddress.Bytes())
	sha.Write(w.Token.Bytes())
	sha.Write(BigToHash(w.Amount).Bytes())
	sha.Write(w.Trader.Bytes())
	sha.Write(w.Receiver.Bytes())
	sha.Write(BigToHash(w.Nonce).Bytes())
	// sha.Write(BigToHash(w.Fee).Bytes())

	return BytesToHash(sha.Sum(nil))
}

func (w *Withdrawal) Sign(wallet *Wallet) error {
	hash := w.ComputeWithdrawalHash()
	signature, err := wallet.SignHash(hash)
	if err != nil {
		return err
	}

	w.Signature = signature
	return nil
}
