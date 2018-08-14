package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Account corresponds to a single Ethereum address. It contains a list of token balances for that address
type Account struct {
	ID            bson.ObjectId                    `json:"id" bson:"_id"`
	Address       common.Address                   `json:"address" bson:"address"`
	TokenBalances map[common.Address]*TokenBalance `json:"tokenBalances" bson:"tokenBalances"`
	IsBlocked     bool                             `json:"isBlocked" bson:"isBlocked"`
	CreatedAt     time.Time                        `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time                        `json:"updatedAt" bson:"updatedAt"`
}

// TokenBalance holds the Balance, Allowance and the Locked balance values for a single Ethereum token
// Balance, Allowance and Locked Balance are stored as big.Int as they represent uint256 values
type TokenBalance struct {
	TokenID       bson.ObjectId  `json:"tokenId" bson:"tokenId"`
	TokenAddress  common.Address `json:"tokenAddress" bson:"tokenAddress"`
	TokenSymbol   string         `json:"tokenSymbol" bson:"tokenSymbol"`
	Balance       *big.Int       `json:"amount" bson:"amount"`
	Allowance     *big.Int       `json:"allowance" bson:"allowance"`
	LockedBalance *big.Int       `json:"lockedAmount" bson:"lockedAmount"`
}

// AccountRecord corresponds to what is stored in the DB. big.Ints are encoded as strings
type AccountRecord struct {
	ID            bson.ObjectId                 `json:"id" bson:"_id"`
	Address       string                        `json:"address" bson:"address"`
	TokenBalances map[string]TokenBalanceRecord `json:"tokenBalances" bson:"tokenBalances"`
	IsBlocked     bool                          `json:"isBlocked" bson:"isBlocked"`
	CreatedAt     time.Time                     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time                     `json:"updatedAt" bson:"updatedAt"`
}

// TokenBalanceRecord corresponds to a TokenBalance struct that is stored in the DB. big.Ints are encoded as strings
type TokenBalanceRecord struct {
	TokenID       bson.ObjectId `json:"tokenId" bson:"tokenId"`
	TokenAddress  string        `json:"tokenAddress" bson:"tokenAddress"`
	TokenSymbol   string        `json:"tokenSymbol" bson:"tokenSymbol"`
	IsBlocked     bool          `json:"isBlocked" bson:"isBlocked"`
	Balance       string        `json:"balance" bson:"balance"`
	Allowance     string        `json:"allowance" bson:"allowance"`
	LockedBalance string        `json:"lockedBalance" bson:"lockedBalance"`
}

// GetBSON implements bson.Getter
func (a *Account) GetBSON() (interface{}, error) {
	tokenBalances := make(map[string]TokenBalanceRecord)

	for key, value := range a.TokenBalances {
		tokenBalances[key.Hex()] = TokenBalanceRecord{
			TokenID:       value.TokenID,
			TokenAddress:  value.TokenAddress.Hex(),
			TokenSymbol:   value.TokenSymbol,
			Balance:       value.Balance.String(),
			Allowance:     value.Allowance.String(),
			LockedBalance: value.LockedBalance.String(),
		}
	}

	return AccountRecord{
		ID:            a.ID,
		Address:       a.Address.Hex(),
		TokenBalances: tokenBalances,
	}, nil
}

// SetBSON implemenets bson.Setter
func (a *Account) SetBSON(raw bson.Raw) error {
	decoded := &AccountRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	a.TokenBalances = make(map[common.Address]*TokenBalance)
	for key, value := range decoded.TokenBalances {

		balance := new(big.Int)
		balance, _ = balance.SetString(value.Balance, 10)
		allowance := new(big.Int)
		allowance, _ = allowance.SetString(value.Allowance, 10)
		lockedBalance := new(big.Int)
		lockedBalance, _ = lockedBalance.SetString(value.LockedBalance, 10)

		a.TokenBalances[common.HexToAddress(key)] = &TokenBalance{
			TokenID:       value.TokenID,
			TokenAddress:  common.HexToAddress(value.TokenAddress),
			TokenSymbol:   value.TokenSymbol,
			Balance:       balance,
			Allowance:     allowance,
			LockedBalance: lockedBalance,
		}
	}

	a.ID = decoded.ID
	a.Address = common.HexToAddress(decoded.Address)
	a.IsBlocked = decoded.IsBlocked
	a.CreatedAt = decoded.CreatedAt
	a.UpdatedAt = decoded.UpdatedAt

	return nil
}

// Validate enforces the account model
func (a *Account) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address, validation.Required),
	)
}
