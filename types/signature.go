package types

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signature struct
type Signature struct {
	V byte
	R common.Hash
	S common.Hash
}

// NewSignature function decodes []byte to Signature type
func NewSignature(b []byte) (*Signature, error) {
	if len(b) != 64 {
		return nil, errors.New("Signature length should be 64 bytes")
	}

	return &Signature{
		R: common.BytesToHash(b[0:32]),
		S: common.BytesToHash(b[32:65]),
		V: b[64] + 27,
	}, nil
}

// MarshalSignature marshals the signature struct to []byte
func (s *Signature) MarshalSignature() ([]byte, error) {
	sigBytes1 := s.R.Bytes()
	sigBytes2 := s.S.Bytes()
	sigBytes3 := s.V - 27

	sigBytes := append([]byte{}, sigBytes1...)
	sigBytes = append(sigBytes, sigBytes2...)
	sigBytes = append(sigBytes, sigBytes3)

	return sigBytes, nil
}

// Verify returns the address that corresponds to the given signature and signed message
func (s *Signature) Verify(hash common.Hash) (common.Address, error) {

	hashBytes := hash.Bytes()
	sigBytes, err := s.MarshalSignature()
	if err != nil {
		return common.Address{}, err
	}

	pubKey, err := crypto.SigToPub(hashBytes, sigBytes)
	if err != nil {
		return common.Address{}, err
	}
	address := crypto.PubkeyToAddress(*pubKey)
	return address, nil
}

// Sign calculates the EDCSA signature corresponding of a hashed message from a given private key
func Sign(hash common.Hash, privKey *ecdsa.PrivateKey) (*Signature, error) {
	sigBytes, err := crypto.Sign(hash.Bytes(), privKey)
	if err != nil {
		return &Signature{}, err
	}

	sig := &Signature{
		R: common.BytesToHash(sigBytes[0:32]),
		S: common.BytesToHash(sigBytes[32:64]),
		V: sigBytes[64] + 27,
	}

	return sig, nil
}

// SignHash also calculates the EDCSA signature of a message but adds an "Ethereum Signed Message" prefix
// https://github.com/ethereum/EIPs/issues/191
func SignHash(hash common.Hash, privKey *ecdsa.PrivateKey) (*Signature, error) {
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		hash.Bytes(),
	)

	sigBytes, err := crypto.Sign(message, privKey)
	if err != nil {
		return &Signature{}, err
	}

	sig := &Signature{
		R: common.BytesToHash(sigBytes[0:32]),
		S: common.BytesToHash(sigBytes[32:64]),
		V: sigBytes[64] + 27,
	}

	return sig, nil
}
