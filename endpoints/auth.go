package endpoints

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/common"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/auth"
)

type Credential struct {
	Address common.Address `json:"address"`
}

// Auth function is responsible for sending the auth JWT token
func Auth(accountService interfaces.AccountService, signingKey string) routing.Handler {
	return func(c *routing.Context) error {
		var credential Credential
		if err := c.Read(&credential); err != nil {
			return errors.Unauthorized(err.Error())
		}

		account := authenticate(credential, accountService)
		if account == nil {
			return errors.Unauthorized("invalid credential")
		}

		token, err := auth.NewJWT(jwt.MapClaims{
			"address":   account.Address,
			"isBlocked": account.IsBlocked,
			"exp":       time.Now().Add(time.Hour * 72).Unix(),
		}, signingKey)
		if err != nil {
			return errors.Unauthorized(err.Error())
		}

		return c.Write(map[string]string{
			"token": token,
		})
	}
}

func authenticate(c Credential, accountService interfaces.AccountService) (acnt *types.Account) {
	acnt, _ = accountService.Validate(c.Address)
	return
}

// JWTHandler handles jwt middleware, sets userAddress field in requestscope
func JWTHandler(c *routing.Context, j *jwt.Token) error {
	userAddress := j.Claims.(jwt.MapClaims)["address"].(string)
	app.GetRequestScope(c).SetUserAddress(common.HexToAddress(userAddress))
	return nil
}
