package types

type OperatorMessage struct {
	MessageType string
	Order       *Order
	Trade       *Trade
	ErrID       int
}

type PendingTradeMessage struct {
	Order *Order
	Trade *Trade
}
