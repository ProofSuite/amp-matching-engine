package types

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/globalsign/mgo/bson"
)

// Wallet holds both the address and the private key of an ethereum account
type Wallet struct {
	ID         bson.ObjectId
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
	Admin      bool
	Operator   bool
}

// NewWallet returns a new wallet object corresponding to a random private key
func NewWallet() *Wallet {
	privateKey, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Wallet{
		Address:    address,
		PrivateKey: privateKey,
	}
}

// NewWalletFromPrivateKey returns a new wallet object corresponding
// to a given private key
func NewWalletFromPrivateKey(key string) *Wallet {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		logger.Error(err)
	}

	return &Wallet{
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}
}

// GetAddress returns the wallet address
func (w *Wallet) GetAddress() string {
	return w.Address.Hex()
}

// GetPrivateKey returns the wallet private key
func (w *Wallet) GetPrivateKey() string {
	return hex.EncodeToString(w.PrivateKey.D.Bytes())
}

func (w *Wallet) Validate() error {
	return nil
}

type WalletRecord struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Address    string        `json:"address" bson:"address"`
	PrivateKey string        `json:"privateKey" bson:"privateKey"`
	Admin      bool          `json:"admin" bson:"admin"`
	Operator   bool          `json:"operator" bson:"operator"`
}

func (w *Wallet) GetBSON() (interface{}, error) {
	return WalletRecord{
		ID:         w.ID,
		Address:    w.Address.Hex(),
		PrivateKey: hex.EncodeToString(w.PrivateKey.D.Bytes()),
		Admin:      w.Admin,
	}, nil
}

func (w *Wallet) SetBSON(raw bson.Raw) error {
	decoded := &WalletRecord{}
	err := raw.Unmarshal(decoded)
	if err != nil {
		logger.Error(err)
		return err
	}

	w.ID = decoded.ID
	w.Address = common.HexToAddress(decoded.Address)
	w.PrivateKey, err = crypto.HexToECDSA(decoded.PrivateKey)
	if err != nil {
		logger.Error(err)
		return err
	}

	w.Admin = decoded.Admin
	w.Operator = decoded.Operator
	return nil
}

// SignHash signs a hashed message with a wallet private key
// and returns it as a Signature object
func (w *Wallet) SignHash(h common.Hash) (*Signature, error) {
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		h.Bytes(),
	)

	sigBytes, err := crypto.Sign(message, w.PrivateKey)
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

func (w *Wallet) SignOrder(o *Order) error {
	hash := o.ComputeHash()
	sig, err := w.SignHash(hash)
	if err != nil {
		return err
	}

	o.Hash = hash
	o.Signature = sig
	return nil
}
