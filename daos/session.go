package daos

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"reflect"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Database struct contains the pointer to mgo.session
// It is a wrapper over mgo to help utilize mgo connection pool
type Database struct {
	Session *mgo.Session
}

// Global instance of Database struct for singleton use
var db *Database
var logger = utils.Logger
var defaultTimeout = 10 * time.Second

func InitSession(session *mgo.Session) (*mgo.Session, error) {
	if db == nil {
		if session == nil {
			if app.Config.EnableTLS {
				session = NewTLSSession()
			} else {
				session = NewSession()
			}
		}

		db = &Database{session}
	}

	return db.Session, nil
}

func InitTLSSession() (*mgo.Session, error) {
	session := NewTLSSession()
	db = &Database{session}
	return db.Session, nil
}

func (d *Database) InitDatabase(session *mgo.Session) {
	d.Session = session
}

func NewSession() *mgo.Session {
	dialInfo := &mgo.DialInfo{
		Addrs:   []string{app.Config.MongoURL},
		Timeout: 15 * time.Second,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	return session
}

func NewTLSSession() *mgo.Session {
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo := &mgo.DialInfo{
		Addrs: []string{
			"ampcluster0-shard-00-00-xzynf.mongodb.net:27017",
			"ampcluster0-shard-00-01-xzynf.mongodb.net:27017",
			"ampcluster0-shard-00-02-xzynf.mongodb.net:27017",
		},
		Timeout:  60 * time.Second,
		Database: "admin",
		Username: app.Config.MongoDBUsername,
		Password: app.Config.MongoDBPassword,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	return session
}

func NewSecureTLSSession() *mgo.Session {
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	roots := x509.NewCertPool()
	//doesn't seem like it's the right file to be used
	// ca, err := ioutil.ReadFile(app.Config.TLSCACertFile)
	// if err != nil {
	// 	panic(err)
	// }

	// roots.AppendCertsFromPEM(ca)

	tlsConfig.RootCAs = roots

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{app.Config.MongoURL},
		Username: app.Config.MongoDBUsername,
		Password: app.Config.MongoDBPassword,
		Database: "proofdex",
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		if err != nil {
			logger.Error(err)
		}
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	return session
}

// Create is a wrapper for mgo.Insert function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Create(dbName, collection string, data ...interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Insert(data...)
	return
}

// GetByID is a wrapper for mgo.FindId function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetByID(dbName, collection string, id bson.ObjectId, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).FindId(id).SetMaxTime(defaultTimeout).One(response)
	return
}

// Get is a wrapper for mgo.Find function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Get(dbName, collection string, query interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).SetMaxTime(defaultTimeout).Skip(offset).Limit(limit).All(response)
	return
}

func (d *Database) Query(dbName, collection string, query interface{}, selector interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).SetMaxTime(defaultTimeout).Skip(offset).Limit(limit).Select(selector).All(response)
	return
}

// GetAndSort is a wrapper for mgo.Find function with SORT function in pipeline.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetAndSort(dbName, collection string, query interface{}, sort []string, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).SetMaxTime(defaultTimeout).Sort(sort...).Skip(offset).Limit(limit).All(response)
	return
}

// Update is a wrapper for mgo.Update function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Update(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Update(query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (d *Database) UpdateAll(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).UpdateAll(query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (d *Database) Upsert(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).Upsert(query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (d *Database) FindAndModify(dbName, collection string, query interface{}, change mgo.Change, response interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).Find(query).SetMaxTime(defaultTimeout).Apply(change, response)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate is a wrapper for mgo.Pipe function.
// It is used to make mongo aggregate pipeline queries
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Aggregate(dbName, collection string, query []bson.M, response interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	collation := mgo.Collation{Locale: "en", NumericOrdering: true}
	result := reflect.ValueOf(response).Interface()

	err := sc.DB(dbName).C(collection).Pipe(query).Collation(&collation).All(result)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Remove removes one document matching a certain query
func (d *Database) Remove(dbName, collection string, query []bson.M) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Remove(query)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// RemoveAll removes all the documents from a collection matching a certain query
func (d *Database) RemoveAll(dbName, collection string, query interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).RemoveAll(query)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// DropCollection drops all the documents in a collection
func (d *Database) DropCollection(dbName, collection string) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).DropCollection()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
