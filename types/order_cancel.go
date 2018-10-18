package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

// OrderCancel is a group of params used for canceling an order previously
// sent to the matching engine. The OrderId and OrderHash must correspond to the
// same order. To be valid and be able to be processed by the matching engine,
// the OrderCancel must include a signature by the Maker of the order corresponding
// to the OrderHash.
type OrderCancel struct {
	OrderHash common.Hash `json:"orderHash"`
	Hash      common.Hash `json:"hash"`
	Signature *Signature  `json:"signature"`
}

// NewOrderCancel returns a new empty OrderCancel object
func NewOrderCancel() *OrderCancel {
	return &OrderCancel{
		Hash:      common.Hash{},
		OrderHash: common.Hash{},
		Signature: &Signature{},
	}
}

// MarshalJSON returns the json encoded byte array representing the OrderCancel struct
func (oc *OrderCancel) MarshalJSON() ([]byte, error) {
	orderCancel := map[string]interface{}{
		"orderHash": oc.OrderHash,
		"hash":      oc.Hash,
		"signature": map[string]interface{}{
			"V": oc.Signature.V,
			"R": oc.Signature.R,
			"S": oc.Signature.S,
		},
	}

	return json.Marshal(orderCancel)
}

func (oc *OrderCancel) String() string {
	return fmt.Sprintf("\nOrderCancel:\nOrderHash: %x\nHash: %x\nSignature.V: %x\nSignature.R: %x\nSignature.S: %x\n\n",
		oc.OrderHash, oc.Hash, oc.Signature.V, oc.Signature.R, oc.Signature.S)
}

// UnmarshalJSON creates an OrderCancel object from a json byte string
func (oc *OrderCancel) UnmarshalJSON(b []byte) error {
	parsed := map[string]interface{}{}

	err := json.Unmarshal(b, &parsed)
	if err != nil {
		return err
	}

	if parsed["orderHash"] == nil {
		return errors.New("Order Hash is missing")
	}
	oc.OrderHash = common.HexToHash(parsed["orderHash"].(string))

	if parsed["hash"] == nil {
		return errors.New("Hash is missing")
	}
	oc.Hash = common.HexToHash(parsed["hash"].(string))

	sig := parsed["signature"].(map[string]interface{})
	oc.Signature = &Signature{
		V: byte(sig["V"].(float64)),
		R: common.HexToHash(sig["R"].(string)),
		S: common.HexToHash(sig["S"].(string)),
	}

	return nil
}

// VerifySignature returns a true value if the OrderCancel object signature
// corresponds to the Maker of the given order
func (oc *OrderCancel) VerifySignature(o *Order) (bool, error) {
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		oc.Hash.Bytes(),
	)

	address, err := oc.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != o.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

func (oc *OrderCancel) GetSenderAddress() (common.Address, error) {
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		oc.Hash.Bytes(),
	)

	address, err := oc.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}

// ComputeHash computes the hash of an order cancel message
func (oc *OrderCancel) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(oc.OrderHash.Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// Sign first computes the order cancel hash, then signs and sets the signature
func (oc *OrderCancel) Sign(w *Wallet) error {
	h := oc.ComputeHash()
	sig, err := w.SignHash(h)
	if err != nil {
		return err
	}

	oc.Hash = h
	oc.Signature = sig
	return nil
}
