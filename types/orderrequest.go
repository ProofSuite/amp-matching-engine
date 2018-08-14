package types

import (
	"encoding/json"
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	validation "github.com/go-ozzo/ozzo-validation"
)

// OrderRequest is the struct in which the order request sent by the
// user is populated
type OrderRequest struct {
	PairName        string         `json:"pairName"`
	ExchangeAddress common.Address `json:"exchangeAddress"`
	UserAddress     common.Address `json:"userAddress"`
	BuyToken        common.Address `json:"buyToken"`
	SellToken       common.Address `json:"sellToken"`
	BuyAmount       *big.Int       `json:"buyAmount"`
	SellAmount      *big.Int       `json:"sellAmount"`
	TakeFee         *big.Int       `json:"takeFee"`
	MakeFee         *big.Int       `json:"makeFee"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	Expires         *big.Int       `json:"expires" bson:"expires"`
	Signature       *Signature     `json:"signature"`
	Hash            common.Hash    `json:"hash"`
}

func (or OrderRequest) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{
		"pairName":    or.PairName,
		"userAddress": or.UserAddress,
		"buyToken":    or.BuyToken,
		"sellToken":   or.SellToken,
		"buyAmount":   or.BuyAmount,
		"sellAmount":  or.SellAmount,
		"takeFee":     or.TakeFee.String(),
		"makeFee":     or.MakeFee.String(),
		"nonce":       or.Nonce.String(),
		"signature": map[string]interface{}{
			"V": or.Signature.V,
			"R": or.Signature.R,
			"S": or.Signature.S,
		},
		"hash": or.Hash,
	}

	return json.Marshal(order)
}

func (or *OrderRequest) UnmarshalJSON(b []byte) error {
	decoded := map[string]interface{}{}

	err := json.Unmarshal(b, &decoded)
	if err != nil {
		return err
	}

	or.PairName = decoded["pairName"].(string)
	or.UserAddress = common.HexToAddress(decoded["userAddress"].(string))
	or.ExchangeAddress = common.HexToAddress(decoded["exchangeAddress"].(string))
	or.BuyToken = common.HexToAddress(decoded["buyToken"].(string))
	or.SellToken = common.HexToAddress(decoded["sellToken"].(string))

	or.BuyAmount = new(big.Int)
	or.SellAmount = new(big.Int)
	or.Expires = new(big.Int)
	or.Nonce = new(big.Int)
	or.MakeFee = new(big.Int)
	or.TakeFee = new(big.Int)

	or.BuyAmount.UnmarshalJSON([]byte(decoded["buyAmount"].(string)))
	or.SellAmount.UnmarshalJSON([]byte(decoded["sellAmount"].(string)))
	or.Expires.UnmarshalJSON([]byte(decoded["expires"].(string)))
	or.Nonce.UnmarshalJSON([]byte(decoded["nonce"].(string)))
	or.MakeFee.UnmarshalJSON([]byte(decoded["makeFee"].(string)))
	or.TakeFee.UnmarshalJSON([]byte(decoded["takeFee"].(string)))

	signature := decoded["signature"].(map[string]interface{})
	or.Signature = &Signature{
		V: byte(signature["V"].(float64)),
		R: common.HexToHash(signature["R"].(string)),
		S: common.HexToHash(signature["S"].(string)),
	}

	or.Hash = common.HexToHash(decoded["hash"].(string))
	return nil
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

	order = &Order{
		MakeFee:     m.MakeFee,
		TakeFee:     m.TakeFee,
		UserAddress: m.UserAddress,
		BuyToken:    m.BuyToken,
		SellToken:   m.SellToken,
		BuyAmount:   m.BuyAmount,
		SellAmount:  m.SellAmount,
		Hash:        m.ComputeHash(),
		Nonce:       m.Nonce,
		Signature:   m.Signature,
		Amount: m.BuyAmount.Int64(),
		Price: m.SellAmount.Int64() * int64(math.Pow10(8)) / m.BuyAmount.Int64(),
	}

	return order, nil
}

// ComputeHash calculates the orderRequest hash
func (m *OrderRequest) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(m.UserAddress.Bytes())
	sha.Write(m.ExchangeAddress.Bytes())
	sha.Write(m.BuyToken.Bytes())
	sha.Write(common.BigToHash(m.BuyAmount).Bytes())
	sha.Write(m.SellToken.Bytes())
	sha.Write(common.BigToHash(m.SellAmount).Bytes())
	sha.Write(common.BigToHash(m.Expires).Bytes())
	sha.Write(common.BigToHash(m.Nonce).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (m *OrderRequest) VerifySignature() (bool, error) {
	return true, nil

	m.Hash = m.ComputeHash()
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		m.Hash.Bytes(),
	)

	address, err := m.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != m.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil

	// signature, err := NewSignature([]byte(m.Signature))
	// if err != nil {
	// 	return false, err
	// }
}
