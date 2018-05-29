package dex

import (
	"math/big"
)

// Nonce alias type represents the number of transactions of an Ethereum account
type Nonce *big.Int

// TokenAmount alias represents the number of tokens (for a certain Token in a certain Address)
type TokenAmount *big.Int

// Fee alias represents a feeTake or a feeMake in an Order object
type Fee *big.Int

// Transactions alias
// type Transaction *types.Transaction
