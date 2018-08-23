package services

import (
	"github.com/Proofsuite/amp-matching-engine/daos"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
	"io/ioutil"
)

var server dbtest.DBServer
var db *mgo.Session

func init() {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	if _, err := daos.InitSession(session); err != nil {
		panic(err)
	}
	db = session
}
