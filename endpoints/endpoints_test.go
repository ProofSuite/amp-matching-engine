package endpoints

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Sirupsen/logrus"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/stretchr/testify/assert"
)

type apiTestCase struct {
	tag      string
	method   string
	url      string
	body     string
	status   int
	response string
}

func newRouter() *routing.Router {
	logger := logrus.New()
	logger.Level = logrus.PanicLevel

	router := routing.New()
	// the test may be started from the home directory or a subdirectory
	err := app.LoadConfig("./config", "../config")
	if err != nil {
		panic(err)
	}
	// connect to the database
	if err := daos.InitSession(); err != nil {
		panic(err)
	}

	router.Use(
		app.Init(logger),
		content.TypeNegotiator(content.JSON),
	)
	return router
}

func testAPI(router *routing.Router, method, URL, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, URL, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func runAPITests(t *testing.T, router *routing.Router, tests []apiTestCase) {
	for _, test := range tests {
		res := testAPI(router, test.method, test.url, test.body)
		assert.Equal(t, test.status, res.Code, test.tag)
		if test.response != "" {
			// assert.JSONEq(t, test.response, res.Body.String(), test.tag)
			assert.Contains(t, test.response, res.Body.String(), test.tag)
		}
	}
}
