package types

import (
	"math/big"
	"time"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Account corresponds to a single Ethereum address. It contains a list of token balances for that address
type Account struct {
	ID            bson.ObjectId                    `json:"-" bson:"_id"`
	Address       common.Address                   `json:"address" bson:"address"`
	TokenBalances map[common.Address]*TokenBalance `json:"tokenBalances" bson:"tokenBalances"`
	IsBlocked     bool                             `json:"isBlocked" bson:"isBlocked"`
	CreatedAt     time.Time                        `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time                        `json:"updatedAt" bson:"updatedAt"`
}

// TokenBalance holds the Balance, Allowance and the Locked balance values for a single Ethereum token
// Balance, Allowance and Locked Balance are stored as big.Int as they represent uint256 values
type TokenBalance struct {
	Address        common.Address `json:"address" bson:"address"`
	Symbol         string         `json:"symbol" bson:"symbol"`
	Balance        *big.Int       `json:"balance" bson:"balance"`
	Allowance      *big.Int       `json:"allowance" bson:"allowance"`
	PendingBalance *big.Int       `json:"pendingBalance" bson:"pendingBalance"`
	LockedBalance  *big.Int       `json:"lockedBalance" bson:"lockedBalance"`
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
	Address        string `json:"address" bson:"address"`
	Symbol         string `json:"symbol" bson:"symbol"`
	Balance        string `json:"balance" bson:"balance"`
	Allowance      string `json:"allowance" bson:"allowance"`
	PendingBalance string `json:"pendingBalance" base:"pendingBalance"`
	LockedBalance  string `json:"lockedBalance" bson:"lockedBalance"`
}

// GetBSON implements bson.Getter
func (a *Account) GetBSON() (interface{}, error) {
	ar := AccountRecord{
		IsBlocked: a.IsBlocked,
		Address:   a.Address.Hex(),
	}

	tokenBalances := make(map[string]TokenBalanceRecord)

	for key, value := range a.TokenBalances {
		tokenBalances[key.Hex()] = TokenBalanceRecord{
			Address:        value.Address.Hex(),
			Symbol:         value.Symbol,
			Balance:        value.Balance.String(),
			Allowance:      value.Allowance.String(),
			LockedBalance:  value.LockedBalance.String(),
			PendingBalance: value.PendingBalance.String(),
		}
	}

	ar.TokenBalances = tokenBalances

	if a.ID.Hex() == "" {
		ar.ID = bson.NewObjectId()
	} else {
		ar.ID = a.ID
	}

	return ar, nil
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
		pendingBalance := new(big.Int)
		pendingBalance, _ = pendingBalance.SetString(value.PendingBalance, 10)

		a.TokenBalances[common.HexToAddress(key)] = &TokenBalance{
			Address:        common.HexToAddress(value.Address),
			Symbol:         value.Symbol,
			Balance:        balance,
			Allowance:      allowance,
			LockedBalance:  lockedBalance,
			PendingBalance: pendingBalance,
		}
	}

	a.Address = common.HexToAddress(decoded.Address)
	a.ID = decoded.ID
	a.IsBlocked = decoded.IsBlocked
	a.CreatedAt = decoded.CreatedAt
	a.UpdatedAt = decoded.UpdatedAt

	return nil
}

// JSON Marshal/Unmarshal interface

// MarshalJSON implements the json.Marshal interface
func (a *Account) MarshalJSON() ([]byte, error) {
	account := map[string]interface{}{
		"id":        a.ID,
		"address":   a.Address,
		"isBlocked": a.IsBlocked,
		"createdAt": a.CreatedAt.String(),
		"updatedAt": a.UpdatedAt.String(),
	}

	tokenBalance := make(map[string]interface{})

	for address, balance := range a.TokenBalances {
		tokenBalance[address.Hex()] = map[string]interface{}{
			"address":        balance.Address.Hex(),
			"symbol":         balance.Symbol,
			"balance":        balance.Balance.String(),
			"allowance":      balance.Allowance.String(),
			"lockedBalance":  balance.LockedBalance.String(),
			"pendingBalance": balance.PendingBalance.String(),
		}
	}

	account["tokenBalances"] = tokenBalance
	return json.Marshal(account)
}

func (a *Account) UnmarshalJSON(b []byte) error {
	account := map[string]interface{}{}
	err := json.Unmarshal(b, &account)
	if err != nil {
		return err
	}

	if account["id"] != nil && bson.IsObjectIdHex(account["id"].(string)) {
		a.ID = bson.ObjectIdHex(account["id"].(string))
	}

	if account["address"] != nil {
		a.Address = common.HexToAddress(account["address"].(string))
	}

	if account["tokenBalances"] != nil {
		tokenBalances := account["tokenBalances"].(map[string]interface{})
		a.TokenBalances = make(map[common.Address]*TokenBalance)
		for address, balance := range tokenBalances {
			if !common.IsHexAddress(address) {
				continue
			}

			tokenBalance := balance.(map[string]interface{})
			tb := &TokenBalance{}

			if tokenBalance["address"] != nil && common.IsHexAddress(tokenBalance["address"].(string)) {
				tb.Address = common.HexToAddress(tokenBalance["address"].(string))
			}

			if tokenBalance["symbol"] != nil {
				tb.Symbol = tokenBalance["symbol"].(string)
			}

			tb.Balance = new(big.Int)
			tb.Allowance = new(big.Int)
			tb.LockedBalance = new(big.Int)
			tb.PendingBalance = new(big.Int)

			if tokenBalance["balance"] != nil {
				tb.Balance.UnmarshalJSON([]byte(tokenBalance["balance"].(string)))
			}

			if tokenBalance["allowance"] != nil {
				tb.Allowance.UnmarshalJSON([]byte(tokenBalance["allowance"].(string)))
			}

			if tokenBalance["lockedBalance"] != nil {
				tb.LockedBalance.UnmarshalJSON([]byte(tokenBalance["lockedBalance"].(string)))
			}

			if tokenBalance["pendingBalance"] != nil {
				tb.PendingBalance.UnmarshalJSON([]byte(tokenBalance["pendingBalance"].(string)))
			}

			a.TokenBalances[common.HexToAddress(address)] = tb
		}
	}

	return nil
}

// Validate enforces the account model
func (a Account) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address, validation.Required),
	)
}

type AccountBSONUpdate struct {
	*Account
}

func (a *AccountBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	tokenBalances := make(map[string]TokenBalanceRecord)

	//TODO validate this. All the fields have to be set
	for key, value := range a.TokenBalances {
		tokenBalances[key.Hex()] = TokenBalanceRecord{
			Address:        value.Address.Hex(),
			Symbol:         value.Symbol,
			Balance:        value.Balance.String(),
			Allowance:      value.Allowance.String(),
			LockedBalance:  value.LockedBalance.String(),
			PendingBalance: value.PendingBalance.String(),
		}
	}

	set := bson.M{
		"updatedAt": now,
		"address":   a.Address,
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}
