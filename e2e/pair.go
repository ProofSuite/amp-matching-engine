package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/stretchr/testify/assert"
)

func testPair(t *testing.T, tokens []types.Token) []types.Pair {
	fmt.Printf("\n=== Starting Pair test ===\n")
	router := buildRouter()
	listPairs := make([]types.Pair, 0)
	neededPair := types.Pair{
		Name:             strings.ToUpper(tokens[0].Symbol + "-" + tokens[1].Symbol),
		BuyToken:         tokens[1].ID,
		BuyTokenAddress:  tokens[1].ContractAddress,
		BuyTokenSymbol:   tokens[1].Symbol,
		SellToken:        tokens[0].ID,
		SellTokenAddress: tokens[0].ContractAddress,
		SellTokenSymbol:  tokens[0].Symbol,
		Active:           true,
		MakerFee:         0,
		TakerFee:         0,
	}

	// create pair test
	res := testAPI(router, "POST", "/pairs", `{"sellTokenAddress":"`+tokens[0].ContractAddress+`", "buyTokenAddress":"`+tokens[1].ContractAddress+`", "active":true}`)
	assert.Equal(t, http.StatusOK, res.Code, "t1 - create pair")
	var resp types.Pair
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if comparePair(t, resp, neededPair) {
		fmt.Println("PASS  't1 - create pair'")
	} else {
		fmt.Println("FAIL  't1 - create pair'")
	}

	listPairs = append(listPairs, neededPair)

	// Duplicate pair test
	res = testAPI(router, "POST", "/pairs", `{"sellTokenAddress":"`+tokens[0].ContractAddress+`", "buyTokenAddress":"`+tokens[1].ContractAddress+`"}`)

	if assert.Equal(t, 401, res.Code, "t2 - create duplicate pair") {
		fmt.Println("PASS  't2 - create duplicate pair'")
	} else {
		fmt.Println("FAIL  't2 - create duplicate pair'")
	}

	// fetch pair detail test
	res = testAPI(router, "GET", "/pairs/"+tokens[1].ContractAddress+"/"+tokens[0].ContractAddress, "")
	assert.Equal(t, http.StatusOK, res.Code, "t2 - fetch pair")
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if comparePair(t, resp, neededPair) {
		fmt.Println("PASS  't2 - fetch pair'")
	} else {
		fmt.Println("FAIL  't2 - fetch pair'")
	}

	// fetch pairs list
	res = testAPI(router, "GET", "/pairs", "")
	assert.Equal(t, http.StatusOK, res.Code, "t3 - fetch pair list")
	var respArray []types.Pair
	if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
		fmt.Printf("%v", err)
	}

	if assert.Lenf(t, respArray, len(listPairs), fmt.Sprintf("Expected Length: %v Got length: %v", len(listPairs), len(respArray))) {
		rb := true
		for i, r := range respArray {
			if rb = rb && comparePair(t, r, listPairs[i]); !rb {
				fmt.Println("FAIL  't3 - fetch pair list'")
				break
			}
		}
		if rb {
			fmt.Println("PASS  't3 - fetch pair list'")
		}
	} else {
		fmt.Println("FAIL  't3 - fetch pair list'")
	}
	return listPairs
}

func comparePair(t *testing.T, actual, expected types.Pair, msgs ...string) bool {
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	response := true
	response = response && assert.Equalf(t, actual.Name, expected.Name, fmt.Sprintf("Pair Name doesn't match. Expected: %v , Got: %v", expected.Name, actual.Name))

	response = response && assert.Equalf(t, actual.BuyToken.Hex(), expected.BuyToken.Hex(), fmt.Sprintf("Pair BuyToken ID doesn't match. Expected: %v , Got: %v", expected.BuyToken, actual.BuyToken))
	response = response && assert.Equalf(t, actual.SellToken.Hex(), expected.SellToken.Hex(), fmt.Sprintf("Pair SellToken ID doesn't match. Expected: %v , Got: %v", expected.SellToken.Hex(), actual.SellToken.Hex()))

	response = response && assert.Equalf(t, actual.BuyTokenAddress, expected.BuyTokenAddress, fmt.Sprintf("Pair BuyToken Address doesn't match. Expected: %v , Got: %v", expected.BuyTokenAddress, actual.BuyTokenAddress))
	response = response && assert.Equalf(t, actual.SellTokenAddress, expected.SellTokenAddress, fmt.Sprintf("Pair SellToken Address doesn't match. Expected: %v , Got: %v", expected.SellTokenAddress, actual.SellTokenAddress))

	response = response && assert.Equalf(t, actual.BuyTokenSymbol, expected.BuyTokenSymbol, fmt.Sprintf("Pair BuyTokenSymbol doesn't match. Expected: %v , Got: %v", expected.BuyTokenSymbol, actual.BuyTokenSymbol))
	response = response && assert.Equalf(t, actual.SellTokenSymbol, expected.SellTokenSymbol, fmt.Sprintf("Pair SellTokenSymbol doesn't match. Expected: %v , Got: %v", expected.SellTokenSymbol, actual.SellTokenSymbol))

	response = response && assert.Equalf(t, actual.Active, expected.Active, fmt.Sprintf("Pair Active doesn't match. Expected: %v , Got: %v", expected.Active, actual.Active))
	response = response && assert.Equalf(t, actual.MakerFee, expected.MakerFee, fmt.Sprintf("Pair MakerFee doesn't match. Expected: %v , Got: %v", expected.MakerFee, actual.MakerFee))
	response = response && assert.Equalf(t, actual.TakerFee, expected.TakerFee, fmt.Sprintf("Pair TakerFee doesn't match. Expected: %v , Got: %v", expected.TakerFee, actual.TakerFee))

	return response
}
