package types

import (
	"math"

	"github.com/Proofsuite/amp-matching-engine/app"
	validation "github.com/go-ozzo/ozzo-validation"
)

type OrderRequest struct {
	Type   int     `json:"type" bson:"type"`
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
	Fee    float64 `json:"fee"`
	// Signature        string  `json:"signature"`
	PairID      string `json:"pairID"`
	PairName    string `json:"pairName"`
	UserAddress string `json:"userAddress"`
}

// Validate validates the OrderRequest fields.
func (m OrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Type, validation.Required, validation.In(1, 2)),
		validation.Field(&m.Amount, validation.Required),
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.UserAddress, validation.Required),
		// validation.Field(&m.Signature, validation.Required),
		// validation.Field(&m.PairID, validation.Required, validation.NewStringRule(bson.IsObjectIdHex, "Invalid pair id")),
		validation.Field(&m.PairName, validation.Required),
	)
}

// ToOrder converts the OrderRequest to Order
func (m *OrderRequest) ToOrder() (order *Order, err error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	order = &Order{
		Type:        OrderType(m.Type),
		Amount:      int64(m.Amount * math.Pow10(8)),
		Price:       int64(m.Price * math.Pow10(8)),
		Fee:         int64(m.Amount * m.Price * (app.Config.TakeFee / 100) * math.Pow10(8)), // amt*price + amt*price*takeFee/100
		PairName:    m.PairName,
		UserAddress: m.UserAddress,

		AmountBuy:  int64(m.Amount * math.Pow10(8)),
		AmountSell: int64(m.Amount * m.Price * math.Pow10(8)),
	}
	return
}
