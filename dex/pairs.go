package dex

import (
	"errors"
	"fmt"

	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Token struct {
	Symbol  string  `json:"symbol"`
	Address Address `json:"address"`
}

type Tokens map[string]Token

type TokenPair struct {
	QuoteToken Token `json:"quoteToken"`
	BaseToken  Token `json:"baseToken"`
	ID         Hash  `json:"id"`
}

type TokenPairs map[string]TokenPair

func NewPair(baseToken, quoteToken Token) TokenPair {
	pair := TokenPair{BaseToken: baseToken, QuoteToken: quoteToken}
	pair.ID = pair.ComputeID()
	return pair
}

func (p *TokenPair) ComputeID() Hash {
	sha := sha3.NewKeccak256()

	sha.Write(p.BaseToken.Address.Bytes())
	sha.Write(p.QuoteToken.Address.Bytes())
	return BytesToHash(sha.Sum(nil))
}

func (p *TokenPair) String() string {
	return fmt.Sprintf("\nQuoteToken: %x at %x\nBaseToken: %x at %x\nID: %x\n\n", p.QuoteToken.Symbol, p.QuoteToken.Address, p.BaseToken.Symbol, p.BaseToken.Address, p.ID)
}

// ComputeOrderPrice calculates the (Amount, Price) tuple corresponding
// to the (AmountBuy, TokenBuy, AmountSell, TokenSell) quadruplet
func (p *TokenPair) ComputeOrderPrice(o *Order) error {
	if o.TokenBuy == p.BaseToken.Address {
		o.OrderType = BUY
		o.Amount = o.AmountBuy.Uint64()
		o.Price = o.AmountSell.Uint64() * (1e3) / o.AmountBuy.Uint64()
	} else if o.TokenBuy == p.QuoteToken.Address {
		o.OrderType = SELL
		o.Amount = o.AmountSell.Uint64()
		o.Price = o.AmountBuy.Uint64() * (1e3) / o.AmountSell.Uint64()
	} else {
		return errors.New("\nToken Buy Address should be either the base token address or the quote token address\n\n")
	}

	return nil
}
