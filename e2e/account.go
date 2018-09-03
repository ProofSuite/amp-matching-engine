package e2e

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func testAccount(t *testing.T, tokens []types.Token) types.Account {
	fmt.Printf("\n=== Starting Account test ===\n")
	router := NewRouter()

	expectedAccount := types.Account{
		Address:       common.HexToAddress("0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"),
		IsBlocked:     false,
		TokenBalances: make(map[common.Address]*types.TokenBalance),
	}
	initBalance := big.NewInt(10000000000000000)

	for _, token := range tokens {
		expectedAccount.TokenBalances[token.ContractAddress] = &types.TokenBalance{
			Address:       token.ContractAddress,
			Symbol:        token.Symbol,
			Balance:       initBalance,
			Allowance:     initBalance,
			LockedBalance: big.NewInt(0),
		}
	}
	// create account test
	res := testAPI(router, "POST", "/account", `{"address":"0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"}`)
	assert.Equal(t, http.StatusOK, res.Code, "t1 - create account")
	var resp types.Account
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if compareAccount(t, resp, expectedAccount) {
		fmt.Println("PASS  't1 - create account'")
	} else {
		fmt.Println("FAIL  't1 - create account'")
	}

	expectedAccount = resp
	// Duplicate account test
	res1 := testAPI(router, "POST", "/account", `{"address":"0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"}`)

	if err := json.Unmarshal(res1.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}

	if assert.Equal(t, expectedAccount.ID.Hex(), resp.ID.Hex(), "t2 - create duplicate account") {
		fmt.Println("PASS  't2 - create duplicate account'")
	} else {
		fmt.Println("FAIL  't2 - create duplicate account'")
	}

	// Get Account
	res2 := testAPI(router, "GET", "/account/0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "")

	assert.Equal(t, http.StatusOK, res2.Code, "t3 - fetch account")
	if err := json.Unmarshal(res2.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if compareAccount(t, resp, expectedAccount) {
		fmt.Println("PASS  't3 - create account'")
	} else {
		fmt.Println("FAIL  't3 - create account'")
	}
	return expectedAccount
}

func compareAccount(t *testing.T, actual, expected types.Account, msgs ...string) bool {
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	response := true
	response = response && assert.Equalf(t, actual.Address, expected.Address, fmt.Sprintf("Address doesn't match. Expected: %v , Got: %v", expected.Address, actual.Address))
	response = response && assert.Equalf(t, actual.IsBlocked, expected.IsBlocked, fmt.Sprintf("Address IsBlocked doesn't match. Expected: %v , Got: %v", expected.IsBlocked, actual.IsBlocked))
	response = response && assert.Equalf(t, actual.TokenBalances, expected.TokenBalances, fmt.Sprintf("Balance Tokens doesn't match. Expected: %v , Got: %v", expected.TokenBalances, actual.TokenBalances))

	return response
}
