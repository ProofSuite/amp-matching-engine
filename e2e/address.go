package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/stretchr/testify/assert"
)

func testAddress(t *testing.T, tokens []types.Token) types.UserAddress {
	fmt.Printf("\n=== Starting Address test ===\n")
	router := buildRouter()

	neededAddress := types.UserAddress{
		Address:   "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
		IsBlocked: false,
	}
	// create address test
	res := testAPI(router, "POST", "/address", `{"address":"0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"}`)
	assert.Equal(t, http.StatusOK, res.Code, "t1 - create address")
	var resp types.UserAddress
	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}
	if compareAddress(t, resp, neededAddress) {
		fmt.Println("PASS  't1 - create address'")
	} else {
		fmt.Println("FAIL  't1 - create address'")
	}

	neededAddress = resp
	// Duplicate address test
	res1 := testAPI(router, "POST", "/address", `{"address":"0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"}`)

	if err := json.Unmarshal(res1.Body.Bytes(), &resp); err != nil {
		fmt.Printf("%v", err)
	}

	if assert.Equal(t, neededAddress.ID.Hex(), resp.ID.Hex(), "t2 - create duplicate address") {
		fmt.Println("PASS  't2 - create duplicate address'")
	} else {
		fmt.Println("FAIL  't2 - create duplicate address'")
	}

	return neededAddress
}

func compareAddress(t *testing.T, actual, expected types.UserAddress, msgs ...string) bool {
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	response := true
	response = response && assert.Equalf(t, actual.Address, expected.Address, fmt.Sprintf("Address doesn't match. Expected: %v , Got: %v", expected.Address, actual.Address))

	response = response && assert.Equalf(t, actual.IsBlocked, expected.IsBlocked, fmt.Sprintf("Address IsBlocked doesn't match. Expected: %v , Got: %v", expected.IsBlocked, actual.IsBlocked))

	return response
}
