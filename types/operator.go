package types

import "github.com/ethereum/go-ethereum/common"

type OperatorMessage struct {
	MessageType string
	Matches     *Matches
	ErrID       int
}

type OperatorTxSuccessMessage struct {
	MessageType string
	OrderHashes []common.Hash
	TradeHashes []common.Hash
}

// type OperatorMessage struct {
// 	MessageType string
// 	Order       *Order
// 	Trade       *Trade
// 	ErrID       int
// }

type PendingTradeBatch struct {
	Matches *Matches
}
