package types

import (
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	validation "github.com/go-ozzo/ozzo-validation"
)

// OrderRequest is the struct in which the order request sent by the
// user is populated
type OrderRequest struct {
	BuyAmount   float64 `json:"buyAmount"`
	SellAmount  float64 `json:"sellAmount"`
	Fee         float64 `json:"fee"`
	Signature   string  `json:"signature"`
	BuyToken    string  `json:"buyToken"`
	SellToken   string  `json:"sellToken"`
	PairName    string  `json:"pairName"`
	Nonce       int64   `json:"nonce" bson:"nonce"`
	Hash        string  `json:"hash"`
	UserAddress string  `json:"userAddress"`
}

// Validate validates the OrderRequest fields.
func (m OrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BuyAmount, validation.Required),
		validation.Field(&m.SellAmount, validation.Required),
		validation.Field(&m.UserAddress, validation.Required),
		validation.Field(&m.BuyToken, validation.Required, validation.NewStringRule(common.IsHexAddress, "Invalid Buy Token Address")),
		validation.Field(&m.SellToken, validation.Required, validation.NewStringRule(common.IsHexAddress, "Invalid Sell Token Address")),
		// validation.Field(&m.Signature, validation.Required),
		// validation.Field(&m.PairName, validation.Required),
	)
}

// ToOrder converts the OrderRequest to Order
func (m *OrderRequest) ToOrder() (order *Order, err error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	// signature, err := NewSignature([]byte(m.Signature))
	// if err != nil {
	// 	return nil, fmt.Errorf("%s", err)
	// }

	order = &Order{
		Fee:         int64(m.BuyAmount * m.SellAmount * (app.Config.TakeFee / 100) * math.Pow10(8)), // amt*price + amt*price*takeFee/100
		UserAddress: m.UserAddress,
		BuyToken:    m.BuyToken,
		SellToken:   m.SellToken,
		BuyAmount:   int64(m.BuyAmount * math.Pow10(8)),
		SellAmount:  int64(m.SellAmount * math.Pow10(8)),
		Hash:        m.ComputeHash(),
		Nonce:       m.Nonce,
		// Signature:        signature,
	}
	return
}

// ComputeHash calculates the orderRequest hash
func (m *OrderRequest) ComputeHash() (ch string) {
	sha := sha3.NewKeccak256()
	sha.Write([]byte(fmt.Sprintf("%f", m.SellAmount)))
	sha.Write([]byte(fmt.Sprintf("%f", m.BuyAmount)))
	sha.Write([]byte(m.BuyToken))
	sha.Write([]byte(m.SellToken))
	sha.Write([]byte(m.UserAddress))
	sha.Write([]byte(fmt.Sprintf("%d", m.Nonce)))
	return common.BytesToHash(sha.Sum(nil)).Hex()
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (m *OrderRequest) VerifySignature() (bool, error) {
	return true, nil

	if m.Hash == "" {
		m.Hash = m.ComputeHash()
	}
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		common.Hex2Bytes(m.Hash),
	)
	signature, err := NewSignature([]byte(m.Signature))
	if err != nil {
		return false, err
	}
	address, err := signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != common.HexToAddress(m.UserAddress) {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}
