package types

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()

	if reflect.TypeOf(*wallet) != reflect.TypeOf(Wallet{}) {
		t.Error("Wallet type is not correct")
	}

	address := wallet.GetAddress()
	if addressLength := len(address); addressLength != 42 {
		t.Error("Expected address length to be 40, but got: ", addressLength)
	}

	privateKey := wallet.GetPrivateKey()
	if privateKeyLength := len(privateKey); privateKeyLength != 64 {
		t.Error("Expected private key length to be 64, but got: ", privateKeyLength)
	}
}

func TestNewWalletFromPrivateKey(t *testing.T) {
	key := "7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660"

	wallet := NewWalletFromPrivateKey(key)
	if address := wallet.GetAddress(); address != "0xE8E84ee367BC63ddB38d3D01bCCEF106c194dc47" {
		t.Error("Expected address to equal 0xE8E84ee367BC63ddB38d3D01bCCEF106c194dc47 but got: ", address)
	}
}

func TestBSON(t *testing.T) {
	key := "7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660"
	w := NewWalletFromPrivateKey(key)
	w.ID = bson.NewObjectId()

	data, err := bson.Marshal(w)
	if err != nil {
		t.Error("some error:", err)
	}

	decoded := &Wallet{}
	bson.Unmarshal(data, decoded)

	assert.Equal(
		t,
		w.ID,
		decoded.ID,
		"ID should be encoded and decoded correctly",
	)

	assert.Equal(
		t,
		decoded.Address.Hex(),
		"0xE8E84ee367BC63ddB38d3D01bCCEF106c194dc47",
		"Address should be encoded and decoded correctly",
	)

	assert.Equal(
		t,
		hex.EncodeToString(decoded.PrivateKey.D.Bytes()),
		"7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660",
		"Private key should be encoded and decoded correctly",
	)
}
