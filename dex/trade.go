package dex

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

// Trade struct holds arguments corresponding to a "Taker Order"
// To be valid an accept by the matching engine (and ultimately the exchange smart-contract),
// the trade signature must be made from the trader Maker account
type Trade struct {
	OrderHash  Hash       `json:"orderHash"`
	Amount     *big.Int   `json:"amount"`
	TradeNonce *big.Int   `json:"tradeNonce"`
	Taker      Address    `json:"taker"`
	Signature  *Signature `json:"signature"`
	Hash       Hash       `json:"hash"`
	PairID     Hash       `json:"pairID"`
	events     chan *Event
}

// NewTrade returns a new unsigned trade corresponding to an Order, amount and taker address
func NewTrade(o *Order, amount *big.Int, taker Address) *Trade {
	return &Trade{
		OrderHash:  o.Hash,
		PairID:     o.PairID,
		Amount:     amount,
		TradeNonce: big.NewInt(0),
		Taker:      taker,
		Signature:  &Signature{},
		Hash:       Hash{},
	}
}

// String return the standard trade format string
func (t *Trade) String() string {
	return fmt.Sprintf("\nTrade:\nOrderHash: %v\nAmount: %v\nTradeNonce: %v\nTaker: %v\nHash: %v\nPairID: %v\n\n",
		t.OrderHash.String(),
		t.Amount,
		t.TradeNonce,
		t.Taker.String(),
		t.Hash.String(),
		t.PairID.String(),
	)
}

// ComputeTradeHash returns hashes the trade
//
// The OrderHash, Aounot, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *Trade) ComputeHash() Hash {
	sha := sha3.NewKeccak256()

	sha.Write(t.OrderHash.Bytes())
	sha.Write(BigToHash(t.Amount).Bytes())
	sha.Write(t.Taker.Bytes())
	sha.Write(BigToHash(t.TradeNonce).Bytes())
	return BytesToHash(sha.Sum(nil))
}

// Sign calculates ands sets the trade hash and signature with the
// given wallet
func (t *Trade) Sign(w *Wallet) error {
	hash := t.ComputeHash()
	signature, err := w.SignHash(hash)
	if err != nil {
		return err
	}

	t.Hash = hash
	t.Signature = signature
	return nil
}

// Valid verifies that all the fields of a struct are set and
// not null
func (t *Trade) Validate() error {
	if t.OrderHash.String() == "" {
		return errors.New("Order Hash missing")
	}

	if t.Hash.String() == "" {
		return errors.New("Trade Hash missing")
	}

	if t.Amount.Sign() == 0 {
		return errors.New("Amount missing or amount null")
	}

	if t.Taker.String() == "" {
		return errors.New("Taker address is not set")
	}

	if t.Signature == nil {
		return errors.New("Signature is not set")
	}

	return nil
}

// VerifySignature verifies that the trade is correct and corresponds
// to the trade Taker address
func (t *Trade) VerifySignature() (bool, error) {
	address, err := t.Signature.Verify(t.Hash)
	if err != nil {
		return false, err
	}

	if address != t.Taker {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *Trade) MarshalJSON() ([]byte, error) {

	trade := map[string]interface{}{
		"orderHash":  t.OrderHash,
		"amount":     t.Amount.String(),
		"tradeNonce": t.TradeNonce.String(),
		"taker":      t.Taker,
		"pairID":     t.PairID,
		"signature": map[string]interface{}{
			"V": t.Signature.V,
			"R": t.Signature.R,
			"S": t.Signature.S,
		},
		"hash": t.Hash,
	}
	return json.Marshal(trade)
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *Trade) UnmarshalJSON(b []byte) error {
	trade := map[string]interface{}{}

	err := json.Unmarshal(b, &trade)
	if err != nil {
		return err
	}

	if trade["orderHash"] == nil {
		return errors.New("Order Hash is not set")
	}
	t.OrderHash = HexToHash(trade["orderHash"].(string))

	log.Printf("Pair ID is equal to %v", err)

	if trade["pairID"] == nil {
		return errors.New("Pair ID is not set")
	}
	t.PairID = HexToHash(trade["pairID"].(string))

	if trade["hash"] == nil {
		return errors.New("Hash is not set")
	}
	t.Hash = HexToHash(trade["hash"].(string))

	t.Amount = new(big.Int)
	t.Amount.UnmarshalJSON([]byte(trade["amount"].(string)))
	t.TradeNonce = new(big.Int)
	t.TradeNonce.UnmarshalJSON([]byte(trade["amount"].(string)))
	t.Taker = HexToAddress(trade["taker"].(string))
	t.Taker.UnmarshalJSON([]byte(trade["taker"].(string)))

	signature := trade["signature"].(map[string]interface{})
	t.Signature = &Signature{
		V: byte(signature["V"].(float64)),
		R: HexToHash(signature["R"].(string)),
		S: HexToHash(signature["S"].(string)),
	}

	return nil
}

// DecodeTrade takes a payload previously unmarshalled from a JSON byte string
// and decodes it into an Trade object
func (t *Trade) Decode(trade map[string]interface{}) error {
	if trade["orderHash"] == nil {
		return errors.New("Order Hash is not set")
	}
	t.OrderHash = HexToHash(trade["orderHash"].(string))

	if trade["pairID"] == nil {
		return errors.New("Pair ID is not set")
	}
	t.PairID = HexToHash(trade["pairID"].(string))

	t.Amount = new(big.Int)
	t.Amount.UnmarshalJSON([]byte(trade["amount"].(string)))
	t.TradeNonce = new(big.Int)
	t.TradeNonce.UnmarshalJSON([]byte(trade["amount"].(string)))
	t.Taker = HexToAddress(trade["taker"].(string))
	t.Taker.UnmarshalJSON([]byte(trade["taker"].(string)))

	signature := trade["signature"].(map[string]interface{})
	t.Signature = &Signature{
		V: byte(signature["V"].(float64)),
		R: HexToHash(signature["R"].(string)),
		S: HexToHash(signature["S"].(string)),
	}

	t.Hash = HexToHash(trade["hash"].(string))
	return nil
}

// NewTradeExecutedEvent is called when a blockchain transaction is created with the
// trade as input
func (t *Trade) NewTradeExecutedEvent(tx *types.Transaction) *Event {
	payload := &TradeExecutedPayload{Trade: t, Tx: tx.Hash()}
	return &Event{eventType: TRADE_EXECUTED, payload: payload}
}

// NewTradeTransactionSuccessful is called when the operator receives a trade event meaning that the
// exchange was performed successfully on the chain.
func (t *Trade) NewTradeTxSuccess(o *Order, tx *types.Transaction) *Event {
	p := &TxSuccessPayload{Order: o, Trade: t, Tx: tx.Hash()}
	return &Event{eventType: TRADE_TX_SUCCESS, payload: p}
}

// NewTradeTransactionError is called when the operator receives a error event meaning that the
// transaction was interrupted.
func (t *Trade) NewTradeTxError(o *Order, errId uint8) *Event {
	p := &TxErrorPayload{Order: o, Trade: t, ErrorId: errId}
	return &Event{eventType: TRADE_TX_ERROR, payload: p}
}
