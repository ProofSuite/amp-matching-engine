package types

type RawOrderBook struct {
	PairName string   `json:"pairName"`
	Orders   []*Order `json:"orders"`
}
