package dex

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignHash(t *testing.T) {
	orderHash := common.HexToHash("0x8aa8e6cbe04f63443a71fd43d511883087df205e7b47f479bc616713d13ce0c7")
	privateKey, _ := crypto.HexToECDSA("c325e0261d889e2cfd581be2eef17405e4a872ef7a69ada4f09e7375a08c556b")

	signature, err := SignHash(orderHash, privateKey)
	if err != nil {
		t.Error("Error signing order hash")
	}

	expectedSignature := Signature{
		V: 28,
		R: common.HexToHash("0x38ae2328af4da36d8f63b1b034e6de91f8488c4c0226754b0107b9a88532b148"),
		S: common.HexToHash("0x7932e9f29fecd5d483bae7b07df57a04b7ca61ce6dae029274f0cc0358cc6b7c"),
	}

	if expectedSignature != *signature {
		t.Error("Signatures are not equal")
	}
}

func TestMarshalSignature(t *testing.T) {
	s := Signature{
		V: 28,
		R: common.HexToHash("0x38ae2328af4da36d8f63b1b034e6de91f8488c4c0226754b0107b9a88532b148"),
		S: common.HexToHash("0x7932e9f29fecd5d483bae7b07df57a04b7ca61ce6dae029274f0cc0358cc6b7c"),
	}

	sigBytes, err := s.MarshalSignature()
	if err != nil {
		t.Error("Error marshaling signature")
	}

	expected, _ := hex.DecodeString("38ae2328af4da36d8f63b1b034e6de91f8488c4c0226754b0107b9a88532b1480x7932e9f29fecd5d483bae7b07df57a04b7ca61ce6dae029274f0cc0358cc6b7c")
	// expected := []byte()

	if !reflect.DeepEqual(sigBytes, expected) {
		t.Errorf("Expected signature to be equal to %x but got %x instead", expected, sigBytes)
	}
}

func TestVerify(t *testing.T) {
	privKey, _ := crypto.HexToECDSA("c325e0261d889e2cfd581be2eef17405e4a872ef7a69ada4f09e7375a08c556b")
	pubKey := privKey.PublicKey
	expectedAddr := crypto.PubkeyToAddress(pubKey)
	data := common.HexToHash("0x8aa8e6cbe04f63443a71fd43d511883087df205e7b47f479bc616713d13ce0c7")

	s, err := Sign(data, privKey)
	if err != nil {
		t.Errorf("Error verifying hash: %v", err)
	}

	address, err := s.Verify(data)
	if err != nil {
		t.Errorf("Error verifying hash: %v", err)
	}

	if address != expectedAddr {
		t.Errorf("Expected address to be equal to %v but got %v instead", expectedAddr, address)
	}
}
