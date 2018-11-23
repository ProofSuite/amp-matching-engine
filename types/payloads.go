package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/go-ozzo/ozzo-validation"
)

// NewOrderPayload is the struct in which the order request sent by the
// user is populated
type NewOrderPayload struct {
	PairName        string         `json:"pairName"`
	ExchangeAddress common.Address `json:"exchangeAddress"`
	UserAddress     common.Address `json:"userAddress"`
	BaseToken       common.Address `json:"baseToken"`
	QuoteToken      common.Address `json:"quoteToken"`
	Side            string         `json:"side"`
	Amount          *big.Int       `json:"amount"`
	PricePoint      *big.Int       `json:"pricepoint"`
	TakeFee         *big.Int       `json:"takeFee"`
	MakeFee         *big.Int       `json:"makeFee"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	Signature       *Signature     `json:"signature"`
	Hash            common.Hash    `json:"hash"`
}

func (p NewOrderPayload) MarshalJSON() ([]byte, error) {
	encoded := map[string]interface{}{
		"pairName":        p.PairName,
		"exchangeAddress": p.ExchangeAddress,
		"userAddress":     p.UserAddress,
		"amount":          p.Amount.String(),
		"pricepoint":      p.PricePoint.String(),
		"side":            p.Side,
		"takeFee":         p.TakeFee.String(),
		"makeFee":         p.MakeFee.String(),
		"nonce":           p.Nonce.String(),
		"signature": map[string]interface{}{
			"V": p.Signature.V,
			"R": p.Signature.R,
			"S": p.Signature.S,
		},
		"hash": p.Hash,
	}

	return json.Marshal(encoded)
}

func (p *NewOrderPayload) UnmarshalJSON(b []byte) error {
	decoded := map[string]interface{}{}

	err := json.Unmarshal(b, &decoded)
	if err != nil {
		return err
	}

	if decoded["pairName"] != nil {
		p.PairName = decoded["pairName"].(string)
	}

	if decoded["userAddress"] != nil {
		p.UserAddress = common.HexToAddress(decoded["userAddress"].(string))
	}

	if decoded["exchangeAddress"] != nil {
		p.ExchangeAddress = common.HexToAddress(decoded["exchangeAddress"].(string))
	}

	if decoded["amount"] != nil {
		p.Amount = math.ToBigInt(decoded["amount"].(string))
	}

	if decoded["pricepoint"] != nil {
		p.PricePoint = math.ToBigInt(decoded["pricepoint"].(string))
	}

	if decoded["nonce"] != nil {
		p.Nonce = math.ToBigInt(decoded["nonce"].(string))
	}

	if decoded["makeFee"] != nil {
		p.MakeFee = math.ToBigInt(decoded["makeFee"].(string))
	}

	if decoded["takeFee"] != nil {
		p.TakeFee = math.ToBigInt(decoded["takeFee"].(string))
	}

	if decoded["side"] != nil {
		p.Side = decoded["side"].(string)
	}

	if decoded["signature"] != nil {
		signature := decoded["signature"].(map[string]interface{})
		p.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	if decoded["hash"] != nil {
		p.Hash = common.HexToHash(decoded["hash"].(string))
	}
	return nil
}

// Validate validates the NewOrderPayload fields.
func (p NewOrderPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Amount, validation.Required),
		validation.Field(&p.PricePoint, validation.Required),
		validation.Field(&p.UserAddress, validation.Required),
		validation.Field(&p.BaseToken, validation.Required),
		validation.Field(&p.QuoteToken, validation.Required),
		validation.Field(&p.Side, validation.Required),
		// validation.Field(&m.Signature, validation.Required),
	)
}

// ToOrder converts the NewOrderPayload to Order
func (p *NewOrderPayload) ToOrder() (o *Order, err error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	o = &Order{
		MakeFee:     p.MakeFee,
		TakeFee:     p.TakeFee,
		UserAddress: p.UserAddress,
		BaseToken:   p.BaseToken,
		QuoteToken:  p.QuoteToken,
		Amount:      p.Amount,
		Side:        p.Side,
		PricePoint:  p.PricePoint,
		Hash:        p.ComputeHash(),
		Nonce:       p.Nonce,
		Signature:   p.Signature,
	}

	return o, nil
}

func (p *NewOrderPayload) EncodedSide() *big.Int {
	if p.Side == "BUY" {
		return big.NewInt(0)
	} else {
		return big.NewInt(1)
	}
}

// ComputeHash calculates the orderRequest hash
func (p *NewOrderPayload) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(p.ExchangeAddress.Bytes())
	sha.Write(p.UserAddress.Bytes())
	sha.Write(p.BaseToken.Bytes())
	sha.Write(p.QuoteToken.Bytes())
	sha.Write(common.BigToHash(p.Amount).Bytes())
	sha.Write(common.BigToHash(p.PricePoint).Bytes())
	sha.Write(common.BigToHash(p.EncodedSide()).Bytes())
	sha.Write(common.BigToHash(p.Nonce).Bytes())
	sha.Write(common.BigToHash(p.TakeFee).Bytes())
	sha.Write(common.BigToHash(p.MakeFee).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (p *NewOrderPayload) VerifySignature() (bool, error) {
	p.Hash = p.ComputeHash()
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		p.Hash.Bytes(),
	)

	address, err := p.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != p.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}
