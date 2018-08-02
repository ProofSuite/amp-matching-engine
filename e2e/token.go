package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/stretchr/testify/assert"
)

func testToken(t *testing.T) []types.Token {
	fmt.Printf("\n=== Starting Token test ===\n")
	router := buildRouter()
	listTokens := make([]types.Token, 0)
	dbTokensList := make([]types.Token, 0)

	// create token test
	res := testAPI(router, "POST", "/tokens", `{  "name":"HotPotCoin", "symbol":"HPC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3","active":true,"quote":true}`)
	assert.Equal(t, http.StatusOK, res.Code, "t1 - create token")
	var resp types.Token
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if compareToken(t, resp, types.Token{Name: "HotPotCoin", Symbol: "HPC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3", Active: true, Quote: true}) {
		fmt.Println("PASS  't1 - create token'")
	} else {
		fmt.Println("FAIL  't1 - create token'")
	}

	listTokens = append(listTokens, types.Token{Name: "HotPotCoin", Symbol: "HPC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3", Active: true, Quote: true})
	dbTokensList = append(dbTokensList, resp)

	// Duplicate token test
	res = testAPI(router, "POST", "/tokens", `{  "name":"HotPotCoin", "symbol":"HPC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3","active":true,"quote":true }`)

	if assert.Equal(t, 401, res.Code, "t2 - create duplicate token") {
		fmt.Println("PASS  't2 - create duplicate token'")
	} else {
		fmt.Println("FAIL  't2 - create duplicate token'")
	}

	// create second token test
	res = testAPI(router, "POST", "/tokens", `{  "name":"Aura.Test", "symbol":"AUT", "decimal":18, "contractAddress":"0x2034842261b82651885751fc293bba7ba5398156","active":true }`)
	assert.Equal(t, http.StatusOK, res.Code, "t3 - create second token")
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if compareToken(t, resp, types.Token{Name: "Aura.Test", Symbol: "AUT", Decimal: 18, ContractAddress: "0x2034842261b82651885751fc293bba7ba5398156", Active: true}) {
		fmt.Println("PASS  't3 - create second token'")
	} else {
		fmt.Println("FAIL  't3 - create second token'")
	}

	listTokens = append(listTokens, types.Token{Name: "Aura.Test", Symbol: "AUT", Decimal: 18, ContractAddress: "0x2034842261b82651885751fc293bba7ba5398156", Active: true})
	dbTokensList = append(dbTokensList, resp)

	// fetch token detail test
	res = testAPI(router, "GET", "/tokens/0x1888a8db0b7db59413ce07150b3373972bf818d3", "")
	assert.Equal(t, http.StatusOK, res.Code, "t4 - fetch token")
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if compareToken(t, resp, types.Token{Name: "HotPotCoin", Symbol: "HPC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3", Active: true, Quote: true}) {
		fmt.Println("PASS  't4 - fetch token'")
	} else {
		fmt.Println("FAIL  't4 - fetch token'")
	}

	// fetch tokens list
	res = testAPI(router, "GET", "/tokens", "")
	assert.Equal(t, http.StatusOK, res.Code, "t5 - fetch token list")
	var respArray []types.Token
	if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
		fmt.Printf("%v", err)
	}

	if assert.Lenf(t, respArray, len(listTokens), fmt.Sprintf("Expected Length: %v Got length: %v", len(listTokens), len(respArray))) {
		rb := true
		for i, r := range respArray {
			if rb = rb && compareToken(t, r, listTokens[i]); !rb {
				fmt.Println("FAIL  't5 - fetch token list'")
				break
			}
		}
		if rb {
			fmt.Println("PASS  't5 - fetch token list'")
		}
	} else {
		fmt.Println("FAIL  't5 - fetch token list'")
	}

	return dbTokensList
}

func compareToken(t *testing.T, actual, expected types.Token, msgs ...string) bool {
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	response := true
	response = response && assert.Equalf(t, actual.Symbol, expected.Symbol, fmt.Sprintf("Token Symbol doesn't match. Expected: %v , Got: %v", expected.Symbol, actual.Symbol))
	response = response && assert.Equalf(t, actual.Name, expected.Name, fmt.Sprintf("Token Name doesn't match. Expected: %v , Got: %v", expected.Name, actual.Name))
	response = response && assert.Equalf(t, actual.Decimal, expected.Decimal, fmt.Sprintf("Token Decimal doesn't match. Expected: %v , Got: %v", expected.Decimal, actual.Decimal))
	response = response && assert.Equalf(t, actual.ContractAddress, expected.ContractAddress, fmt.Sprintf("Token ContractAddress doesn't match. Expected: %v , Got: %v", expected.ContractAddress, actual.ContractAddress))
	response = response && assert.Equalf(t, actual.Active, expected.Active, fmt.Sprintf("Token Active doesn't match. Expected: %v , Got: %v", expected.Active, actual.Active))
	response = response && assert.Equalf(t, actual.Quote, expected.Quote, fmt.Sprintf("Token Quote doesn't match. Expected: %v , Got: %v", expected.Quote, actual.Quote))

	return response
}
