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

type OrderRequest struct {
	Type        int     `json:"type" bson:"type"`
	Amount      float64 `json:"amount"`
	Price       float64 `json:"price"`
	Fee         float64 `json:"fee"`
	Signature   string  `json:"signature"`
	TokenBuy    string  `json:"tokenBuy"`
	TokenSell   string  `json:"tokenSell"`
	PairName    string  `json:"pairName"`
	Hash        string  `json:"hash"`
	UserAddress string  `json:"userAddress"`
}

// Validate validates the OrderRequest fields.
func (m OrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Type, validation.Required, validation.In(1, 2)),
		validation.Field(&m.Amount, validation.Required),
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.UserAddress, validation.Required),
		validation.Field(&m.TokenBuy, validation.Required, validation.NewStringRule(common.IsHexAddress, "Invalid Buy Token Address")),
		validation.Field(&m.TokenSell, validation.Required, validation.NewStringRule(common.IsHexAddress, "Invalid Sell Token Address")),
		validation.Field(&m.Signature, validation.Required),
		// validation.Field(&m.PairName, validation.Required),
	)
}

// ToOrder converts the OrderRequest to Order
func (m *OrderRequest) ToOrder() (order *Order, err error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	signature, err := NewSignature([]byte(m.Signature))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	order = &Order{
		Type:             OrderType(m.Type),
		Amount:           int64(m.Amount * math.Pow10(8)),
		Price:            int64(m.Price * math.Pow10(8)),
		Fee:              int64(m.Amount * m.Price * (app.Config.TakeFee / 100) * math.Pow10(8)), // amt*price + amt*price*takeFee/100
		UserAddress:      m.UserAddress,
		BuyTokenAddress:  m.TokenBuy,
		SellTokenAddress: m.TokenSell,
		AmountBuy:        int64(m.Amount * math.Pow10(8)),
		AmountSell:       int64(m.Amount * m.Price * math.Pow10(8)),
		Hash:             m.ComputeHash(),
		Signature:        signature,
	}
	return
}

// ComputeHash calculates the orderRequest hash
func (m *OrderRequest) ComputeHash() (ch string) {
	sha := sha3.NewKeccak256()
	sha.Write([]byte(fmt.Sprintf("%f", m.Price)))
	sha.Write([]byte(fmt.Sprintf("%f", m.Amount)))
	sha.Write([]byte(fmt.Sprintf("%d", m.Type)))
	sha.Write([]byte(m.TokenBuy))
	sha.Write([]byte(m.TokenSell))
	sha.Write([]byte(m.UserAddress))
	return common.BytesToHash(sha.Sum(nil)).Hex()
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (m *OrderRequest) VerifySignature() (bool, error) {

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
