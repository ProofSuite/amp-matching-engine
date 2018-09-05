package app

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-dbx"
)

// RequestScope contains the application-specific information that are carried around in a request.
type RequestScope interface {
	Logger
	// UserID returns the ID of the user for the current request
	UserAddress() common.Address
	// SetUserID sets the ID of the currently authenticated user
	SetUserAddress(address common.Address)
	// RequestID returns the ID of the current request
	RequestID() string
	// Tx returns the currently active database transaction that can be used for DB query purpose
	Tx() *dbx.Tx
	// SetTx sets the database transaction
	SetTx(tx *dbx.Tx)
	// Now returns the timestamp representing the time when the request is being processed
	Now() time.Time
}

type requestScope struct {
	Logger                     // the logger tagged with the current request information
	now         time.Time      // the time when the request is being processed
	requestID   string         // an ID identifying one or multiple correlated HTTP requests
	userAddress common.Address // an ID identifying the current user
	tx          *dbx.Tx        // the currently active transaction
}

func (rs *requestScope) UserAddress() common.Address {
	return rs.userAddress
}

func (rs *requestScope) SetUserAddress(address common.Address) {
	rs.Logger.SetField("UserAddress", address.Hex())
	rs.userAddress = address
}

func (rs *requestScope) RequestID() string {
	return rs.requestID
}

func (rs *requestScope) Tx() *dbx.Tx {
	return rs.tx
}

func (rs *requestScope) SetTx(tx *dbx.Tx) {
	rs.tx = tx
}

func (rs *requestScope) Now() time.Time {
	return rs.now
}

// newRequestScope creates a new RequestScope with the current request information.
func newRequestScope(now time.Time, logger *logrus.Logger, request *http.Request) RequestScope {
	l := NewLogger(logger, logrus.Fields{})
	requestID := request.Header.Get("X-Request-Id")
	if requestID != "" {
		l.SetField("RequestID", requestID)
	}
	return &requestScope{
		Logger:    l,
		now:       now,
		requestID: requestID,
	}
}
