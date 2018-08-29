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

func (p NewOrderPayload) MarshalJSON() ([]byte, error) {
	encoded := map[string]interface{}{
		"pairName":        p.PairName,
		"exchangeAddress": p.ExchangeAddress,
		"userAddress":     p.UserAddress,
		"buyToken":        p.BuyToken,
		"sellToken":       p.SellToken,
		"buyAmount":       p.BuyAmount.String(),
		"sellAmount":      p.SellAmount.String(),
		"takeFee":         p.TakeFee.String(),
		"makeFee":         p.MakeFee.String(),
		"expires":         p.Expires.String(),
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

	if decoded["buyToken"] != nil {
		p.BuyToken = common.HexToAddress(decoded["buyToken"].(string))
	}

	if decoded["sellToken"] != nil {
		p.SellToken = common.HexToAddress(decoded["sellToken"].(string))
	}

	if decoded["buyAmount"] != nil {
		p.BuyAmount = math.ToBigInt(decoded["buyAmount"].(string))
	}

	if decoded["sellAmount"] != nil {
		p.BuyAmount = math.ToBigInt(decoded["sellAmount"].(string))
	}

	if decoded["expires"] != nil {
		p.BuyAmount = math.ToBigInt(decoded["expires"].(string))
	}

	if decoded["nonce"] != nil {
		p.BuyAmount = math.ToBigInt(decoded["nonce"].(string))
	}

	if decoded["makeFee"] != nil {
		p.BuyAmount = math.ToBigInt(decoded["makeFee"].(string))
	}

	if decoded["takeFee"] != nil {
		p.BuyAmount = math.ToBigInt(decoded["takeFee"].(string))
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
		validation.Field(&p.BuyAmount, validation.Required),
		validation.Field(&p.SellAmount, validation.Required),
		validation.Field(&p.UserAddress, validation.Required),
		validation.Field(&p.BuyToken, validation.Required),
		validation.Field(&p.SellToken, validation.Required),
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
		BuyToken:    p.BuyToken,
		SellToken:   p.SellToken,
		BuyAmount:   p.BuyAmount,
		SellAmount:  p.SellAmount,
		Hash:        p.ComputeHash(),
		Nonce:       p.Nonce,
		Expires:     p.Expires,
		Signature:   p.Signature,
	}

	return o, nil
}

// ComputeHash calculates the orderRequest hash
func (p *NewOrderPayload) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(p.UserAddress.Bytes())
	sha.Write(p.ExchangeAddress.Bytes())
	sha.Write(p.BuyToken.Bytes())
	sha.Write(common.BigToHash(p.BuyAmount).Bytes())
	sha.Write(p.SellToken.Bytes())
	sha.Write(common.BigToHash(p.SellAmount).Bytes())
	sha.Write(common.BigToHash(p.Expires).Bytes())
	sha.Write(common.BigToHash(p.Nonce).Bytes())
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
