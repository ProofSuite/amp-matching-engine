package e2e

import (
	"fmt"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/stretchr/testify/assert"
)

func testPair(t *testing.T) {
	fmt.Printf("\n=== Starting Pair test ===\n")
	// router := buildRouter()
	// listPairs := make([]types.Pair, 0)

	// create pair test
	// res := testAPI(router, "POST", "/pairs", `{  "name":"ABC", "symbol":"ABC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3" }`)
	// assert.Equal(t, http.StatusOK, res.Code, "t1 - create pair")
	// var resp types.Pair
	// if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
	// 	fmt.Printf("%v", err)
	// }
	// if comparePair(t, resp, types.Pair{Name: "ABC", Symbol: "ABC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3"}) {
	// 	fmt.Println("PASS  't1 - create pair'")
	// } else {
	// 	fmt.Println("FAIL  't1 - create pair'")
	// }

	// listPairs = append(listPairs, types.Pair{Name: "ABC", Symbol: "ABC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3"})

	// // fetch pair detail test
	// res = testAPI(router, "GET", "/pairs/0x1888a8db0b7db59413ce07150b3373972bf818d3", "")
	// assert.Equal(t, http.StatusOK, res.Code, "t2 - fetch pair")
	// if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
	// 	fmt.Printf("%v", err)
	// }
	// if comparePair(t, resp, types.Pair{Name: "ABC", Symbol: "ABC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3"}) {
	// 	fmt.Println("PASS  't2 - fetch pair'")
	// } else {
	// 	fmt.Println("FAIL  't2 - fetch pair'")
	// }

	// // fetch pairs list
	// res = testAPI(router, "GET", "/pairs", "")
	// assert.Equal(t, http.StatusOK, res.Code, "t3 - fetch pair list")
	// var respArray []types.Pair
	// if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
	// 	fmt.Printf("%v", err)
	// }

	// if assert.Lenf(t, respArray, len(listPairs), fmt.Sprintf("Expected Length: %v Got length: %v", len(listPairs), len(respArray))) {
	// 	rb := true
	// 	for i, r := range respArray {
	// 		if rb = rb && comparePair(t, r, listPairs[i]); !rb {
	// 			fmt.Println("FAIL  't3 - fetch pair list'")
	// 			break
	// 		}
	// 	}
	// 	if rb {
	// 		fmt.Println("PASS  't3 - fetch pair list'")
	// 	}
	// } else {
	// 	fmt.Println("FAIL  't3 - fetch pair list'")
	// }

}

func comparePair(t *testing.T, actual, expected types.Pair, msgs ...string) bool {
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	response := true
	response = response && assert.Equalf(t, actual.Name, expected.Name, fmt.Sprintf("Pair Name doesn't match. Expected: %v , Got: %v", expected.Name, actual.Name))

	return response
}
