package types

// // Balance holds both the address and the private key of an ethereum account
// type Balance struct {
// 	ID        bson.ObjectId           `json:"id" bson:"_id"`
// 	Address   string                  `json:"address" bson:"address"`
// 	Tokens    map[string]TokenBalance `json:"tokens" bson:"tokens"`
// 	CreatedAt time.Time               `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt time.Time               `json:"updatedAt" bson:"updatedAt"`
// }

// // TokenBalance is a subdocument
// // It contains the confirmed amount and locked amount
// // corresponding to a single token (identified using tokenID & token's contract address)
// type TokenBalance struct {
// 	ID      bson.ObjectId `json:"tokenId" bson:"tokenId"`
// 	Address string        `json:"tokenAddress" bson:"tokenAddress"`
// 	Symbol  string        `json:"tokenSymbol" bson:"tokenSymbol"`
// 	Amount       int64         `json:"amount" bson:"amount"`
// 	LockedAmount int64         `json:"lockedAmount" bson:"lockedAmount"`
// }

// // NewBalance returns a new wallet object corresponding to a random private key
// func NewBalance(address string) (w *Balance, err error) {
// 	if !common.IsHexAddress(address) {
// 		return nil, errors.New("Invalid Address")
// 	}
// 	return
// }
