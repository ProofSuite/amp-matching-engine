package endpoints

import (
	"net/http"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/services"
)

func TestToken(t *testing.T) {
	router := newRouter()
	ServeTokenResource(&router.RouteGroup, services.NewTokenService(daos.NewTokenDao()))

	// notFoundError := `{"error_code":"NOT_FOUND", "message":"NOT_FOUND"}`
	// nameRequiredError := `{"error_code":"INVALID_DATA","message":"INVALID_DATA","details":[{"field":"name","error":"cannot be blank"}]}`

	runAPITests(t, router, []apiTestCase{
		{"t1 - create token", "POST", "/tokens/5b3e82607b44576ba8000002", `{ "code":"abc", "name":"HotPotCoin", "symbol":"ABC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818e3" }`, http.StatusOK, `{ "code":"abc", "name":"HotPotCoin", "symbol":"ABC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3" }`},
	})
}
