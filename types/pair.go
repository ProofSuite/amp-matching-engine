package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Pair struct is used to model the pair data in the system and DB
type Pair struct {
	ID                bson.ObjectId  `json:"-" bson:"_id"`
	BaseTokenSymbol   string         `json:"baseTokenSymbol,omitempty" bson:"baseTokenSymbol"`
	BaseTokenAddress  common.Address `json:"baseTokenAddress,omitempty" bson:"baseTokenAddress"`
	BaseTokenDecimal  int            `json:"baseTokenDecimal,omitempty" bson:"baseTokenDecimal"`
	QuoteTokenSymbol  string         `json:"quoteTokenSymbol,omitempty" bson:"quoteTokenSymbol"`
	QuoteTokenAddress common.Address `json:"quoteTokenAddress,omitempty" bson:"quoteTokenAddress"`
	QuoteTokenDecimal int            `json:"quoteTokenDecimal,omitempty" bson:"quoteTokenDecimal"`
	PriceMultiplier   *big.Int       `json:"priceMultiplier,omitempty" bson:"priceMultiplier"`
	Active            bool           `json:"active,omitempty" bson:"active"`
	MakeFee           *big.Int       `json:"makeFee,omitempty" bson:"makeFee"`
	TakeFee           *big.Int       `json:"takeFee,omitempty" bson:"takeFee"`
	CreatedAt         time.Time      `json:"-" bson:"createdAt"`
	UpdatedAt         time.Time      `json:"-" bson:"updatedAt"`
}

type PairAddresses struct {
	Name       string         `json:"name" bson:"name"`
	BaseToken  common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken common.Address `json:"quoteToken" bson:"quoteToken"`
}

type PairAddressesRecord struct {
	Name       string `json:"name" bson:"name"`
	BaseToken  string `json:"baseToken" bson:"baseToken"`
	QuoteToken string `json:"quoteToken" bson:"quoteToken"`
}

type PairRecord struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	BaseTokenSymbol   string    `json:"baseTokenSymbol" bson:"baseTokenSymbol"`
	BaseTokenAddress  string    `json:"baseTokenAddress" bson:"baseTokenAddress"`
	BaseTokenDecimal  int       `json:"baseTokenDecimal" bson:"baseTokenDecimal"`
	QuoteTokenSymbol  string    `json:"quoteTokenSymbol" bson:"quoteTokenSymbol"`
	QuoteTokenAddress string    `json:"quoteTokenAddress" bson:"quoteTokenAddress"`
	QuoteTokenDecimal int       `json:"quoteTokenDecimal" bson:"quoteTokenDecimal"`
	Active            bool      `json:"active" bson:"active"`
	PriceMultiplier   string    `json:"priceMultiplier" bson:"priceMultiplier"`
	MakeFee           string    `json:"makeFee" bson:"makeFee"`
	TakeFee           string    `json:"takeFee" bson:"takeFee"`
	CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (p *Pair) Code() string {
	code := p.BaseTokenSymbol + "/" + p.QuoteTokenSymbol + "::" + p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex()
	return code
}

func (p *Pair) AddressCode() string {
	code := p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex()
	return code
}

func (p *Pair) Name() string {
	name := p.BaseTokenSymbol + "/" + p.QuoteTokenSymbol
	return name
}

func (p *Pair) SetBSON(raw bson.Raw) error {
	decoded := &PairRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	makeFee := big.NewInt(0)
	makeFee, _ = makeFee.SetString(decoded.MakeFee, 10)
	takeFee := big.NewInt(0)
	takeFee, _ = takeFee.SetString(decoded.TakeFee, 10)
	priceMultiplier := big.NewInt(0)
	priceMultiplier, _ = priceMultiplier.SetString(decoded.PriceMultiplier, 10)

	p.ID = decoded.ID
	p.BaseTokenSymbol = decoded.BaseTokenSymbol
	p.BaseTokenAddress = common.HexToAddress(decoded.BaseTokenAddress)
	p.BaseTokenDecimal = decoded.BaseTokenDecimal
	p.QuoteTokenSymbol = decoded.QuoteTokenSymbol
	p.QuoteTokenAddress = common.HexToAddress(decoded.QuoteTokenAddress)
	p.QuoteTokenDecimal = decoded.QuoteTokenDecimal
	p.Active = decoded.Active
	p.PriceMultiplier = priceMultiplier
	p.MakeFee = makeFee
	p.TakeFee = takeFee

	p.CreatedAt = decoded.CreatedAt
	p.UpdatedAt = decoded.UpdatedAt
	return nil
}

func (p *Pair) GetBSON() (interface{}, error) {
	return &PairRecord{
		ID: p.ID,

		BaseTokenSymbol:   p.BaseTokenSymbol,
		BaseTokenAddress:  p.BaseTokenAddress.Hex(),
		BaseTokenDecimal:  p.BaseTokenDecimal,
		QuoteTokenSymbol:  p.QuoteTokenSymbol,
		QuoteTokenAddress: p.QuoteTokenAddress.Hex(),
		QuoteTokenDecimal: p.QuoteTokenDecimal,
		PriceMultiplier:   p.PriceMultiplier.String(),
		Active:            p.Active,
		MakeFee:           p.MakeFee.String(),
		TakeFee:           p.TakeFee.String(),
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}, nil
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (p Pair) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BaseTokenAddress, validation.Required),
		validation.Field(&p.QuoteTokenAddress, validation.Required),
		validation.Field(&p.BaseTokenSymbol, validation.Required),
		validation.Field(&p.QuoteTokenSymbol, validation.Required),
	)
}

// GetOrderBookKeys returns the orderbook price point keys for corresponding pair
// It is used to fetch the orderbook from redis of a pair
func (p *Pair) GetOrderBookKeys() (sell, buy string) {
	return p.GetKVPrefix() + "::SELL", p.GetKVPrefix() + "::BUY"
}

func (p *Pair) GetKVPrefix() string {
	return p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex()
}
