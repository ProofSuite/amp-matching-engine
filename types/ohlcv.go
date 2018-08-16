package types

// Tick is the format in which mongo aggregate pipeline returns data when queried for OHLCV data
type Tick struct {
	ID    TickID `json:"_id,omitempty" bson:"_id"`
	C     int64  `json:"c" bson:"c"`
	Count int64  `json:"count" bson:"count"`
	H     int64  `json:"h" bson:"h"`
	L     int64  `json:"l" bson:"l"`
	O     int64  `json:"o" bson:"o"`
	Ts    int64  `json:"ts" bson:"ts"`
	V     int64  `json:"v" bson:"v"`
}

// TickID is the subdocument for aggregate grouping for OHLCV data
type TickID struct {
	Pair       string `json:"pair" bson:"pair"`
	BaseToken  string `json:"baseToken" bson:"baseToken"`
	QuoteToken string `json:"quoteToken" bson:"quoteToken"`
}

type TickRequest struct {
	Pair     []PairSubDoc `json:"pair"`
	From     int64        `json:"from"`
	To       int64        `json:"to"`
	Duration int64        `json:"duration"`
	Units    string       `json:"units"`
}
