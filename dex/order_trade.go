package dex

type OrderTradePair struct {
	order *Order `json:"order"`
	trade *Trade `json:"trade"`
}
