package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Token struct is used to model the token data in the system and DB
type Token struct {
	ID              bson.ObjectId  `json:"-" bson:"_id"`
	Name            string         `json:"name" bson:"name"`
	Symbol          string         `json:"symbol" bson:"symbol"`
	Image           Image          `json:"image" bson:"image"`
	ContractAddress common.Address `json:"contractAddress" bson:"contractAddress"`
	Decimal         int            `json:"decimal" bson:"decimal"`
	Active          bool           `json:"active" bson:"active"`
	Quote           bool           `json:"quote" bson:"quote"`
	MakeFee         string         `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee         string         `json:"takeFee,omitempty" bson:"makeFee,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// TokenRecord is the struct which is stored in db
type TokenRecord struct {
	ID              bson.ObjectId `json:"-" bson:"_id"`
	Name            string        `json:"name" bson:"name"`
	Symbol          string        `json:"symbol" bson:"symbol"`
	Image           Image         `json:"image" bson:"image"`
	ContractAddress string        `json:"contractAddress" bson:"contractAddress"`
	Decimal         int           `json:"decimal" bson:"decimal"`
	Active          bool          `json:"active" bson:"active"`
	Quote           bool          `json:"quote" bson:"quote"`
	MakeFee         string        `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee         string        `json:"takeFee,omitempty" bson:"takeFee,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Image is a sub document used to store data related to images
type Image struct {
	URL  string                 `json:"url" bson:"url"`
	Meta map[string]interface{} `json:"meta" bson:"meta"`
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (t Token) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Symbol, validation.Required),
		validation.Field(&t.ContractAddress, validation.Required),
		validation.Field(&t.Decimal, validation.Required),
	)
}

// GetBSON implements bson.Getter
func (t *Token) GetBSON() (interface{}, error) {

	return TokenRecord{
		ID:              t.ID,
		Name:            t.Name,
		Symbol:          t.Symbol,
		Image:           t.Image,
		ContractAddress: t.ContractAddress.Hex(),
		Decimal:         t.Decimal,
		Active:          t.Active,
		Quote:           t.Quote,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		MakeFee:         t.MakeFee,
		TakeFee:         t.TakeFee,
	}, nil
}

// SetBSON implemenets bson.Setter
func (t *Token) SetBSON(raw bson.Raw) error {
	decoded := &TokenRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}
	t.ID = decoded.ID
	t.Name = decoded.Name
	t.Symbol = decoded.Symbol
	t.Image = decoded.Image
	if common.IsHexAddress(decoded.ContractAddress) {
		t.ContractAddress = common.HexToAddress(decoded.ContractAddress)
	}
	t.Decimal = decoded.Decimal
	t.Active = decoded.Active
	t.Quote = decoded.Quote
	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	t.MakeFee = decoded.MakeFee
	t.TakeFee = decoded.TakeFee
	return nil
}
