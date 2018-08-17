package e2e

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"net/http"
// 	"testing"

// 	"github.com/Proofsuite/amp-matching-engine/types"
// 	"github.com/stretchr/testify/assert"
// )

// func testBalance(t *testing.T, tokens []types.Token, address types.UserAddress) types.Balance {
// 	fmt.Printf("\n=== Starting Balance test ===\n")
// 	router := buildRouter()

// 	tokenBalance := make(map[string]types.TokenBalance)
// 	for _, t := range tokens {
// 		tokenBalance[t.ContractAddress] = types.TokenBalance{
// 			ID:      t.ID,
// 			Address: t.ContractAddress,
// 			Symbol:  t.Symbol,
// 			Amount:       int64(10000 * math.Pow10(8)),
// 			LockedAmount: 0,
// 		}
// 	}
// 	neededBalance := types.Balance{
// 		Address: address.Address,
// 		Tokens:  tokenBalance,
// 	}
// 	// create balance test
// 	res := testAPI(router, "GET", "/balances/"+address.Address, ``)
// 	assert.Equal(t, http.StatusOK, res.Code, "t1 - fetch balance")
// 	var resp types.Balance
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	if compareBalance(t, resp, neededBalance) {
// 		fmt.Println("PASS  't1 - fetch balance'")
// 	} else {
// 		fmt.Println("FAIL  't1 - fetch balance'")
// 	}
// 	return neededBalance
// }

// func compareBalance(t *testing.T, actual, expected types.Balance, msgs ...string) bool {
// 	for _, msg := range msgs {
// 		fmt.Println(msg)
// 	}
// 	response := true
// 	response = response && assert.Equalf(t, actual.Address, expected.Address, fmt.Sprintf("Address doesn't match. Expected: %v , Got: %v", expected.Address, actual.Address))

// 	response = response && assert.Equalf(t, actual.Tokens, expected.Tokens, fmt.Sprintf("Balance Tokens doesn't match. Expected: %v , Got: %v", expected.Tokens, actual.Tokens))

// 	return response
// }
