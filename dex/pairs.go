package dex

import (
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
	return fmt.Sprintf("\nQuoteToken: %v at %x\nBaseToken: %v at %x\nID: %x\n\n", p.QuoteToken.Symbol, p.QuoteToken.Address, p.BaseToken.Symbol, p.BaseToken.Address, p.ID)
}
