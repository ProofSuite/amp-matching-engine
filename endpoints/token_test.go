package endpoints

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/Proofsuite/amp-matching-engine/daos"
// 	"github.com/Proofsuite/amp-matching-engine/services"
// 	"github.com/Proofsuite/amp-matching-engine/types"
// )

// func TestToken(t *testing.T) {
// 	router := newRouter()
// 	ServeTokenResource(&router.RouteGroup, services.NewTokenService(daos.NewTokenDao()))

// 	// notFoundError := `{"error_code":"NOT_FOUND", "message":"NOT_FOUND"}`
// 	// nameRequiredError := `{"error_code":"INVALID_DATA","message":"INVALID_DATA","details":[{"field":"name","error":"cannot be blank"}]}`

// 	runAPITests(t, router, []apiTestCase{
// 		{"t1 - create token", "POST", "/tokens", `{  "name":"ABC", "symbol":"ABC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3" }`, http.StatusOK, types.Token{Name: "ABC", Symbol: "ABC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3"}, `custom`, compareToken},
// 		{"t1 - fetch tokens", "GET", "/tokens", "", http.StatusOK, types.Token{Name: "ABC", Symbol: "ABC", Decimal: 18, ContractAddress: "0x1888a8db0b7db59413ce07150b3373972bf818d3"}, `custom`, compareToken},
// 	})
// }

// func compareToken(t *testing.T, result, testResult interface{}) {
// 	var actual types.Token
// 	actualAsBytes, _ := json.Marshal(result)
// 	json.Unmarshal(actualAsBytes, &actual)
// 	var expected types.Token
// 	expectedAsBytes, _ := json.Marshal(testResult)
// 	json.Unmarshal(expectedAsBytes, &expected)

// 	assert.Equalf(t, actual.Symbol, expected.Symbol, fmt.Sprintf("Token Symbol doesn't match. Expected: %v , Got: %v", expected.Symbol, actual.Symbol))
// 	assert.Equalf(t, actual.Name, expected.Name, fmt.Sprintf("Token Name doesn't match. Expected: %v , Got: %v", expected.Name, actual.Name))
// 	assert.Equalf(t, actual.Decimal, expected.Decimal, fmt.Sprintf("Token Decimal doesn't match. Expected: %v , Got: %v", expected.Decimal, actual.Decimal))
// 	assert.Equalf(t, actual.ContractAddress, expected.ContractAddress, fmt.Sprintf("Token ContractAddress doesn't match. Expected: %v , Got: %v", expected.ContractAddress, actual.ContractAddress))
// 	assert.Equalf(t, actual.Image, expected.Image, fmt.Sprintf("Token Image doesn't match. Expected: %v , Got: %v", expected.Image, actual.Image))
// 	assert.Equalf(t, actual.Active, expected.Active, fmt.Sprintf("Token Active doesn't match. Expected: %v , Got: %v", expected.Active, actual.Active))
// }
