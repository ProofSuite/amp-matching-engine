package daos

import (
	"fmt"

	"github.com/Proofsuite/amp-matching-engine/app"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Database struct {
	session *mgo.Session
}

var DB *Database

func InitSession() error {
	if DB == nil {
		db, err := mgo.Dial(app.Config.DSN)
		if err != nil {
			return err
		}
		DB = &Database{db}
	}
	return nil
}

func (d *Database) Create(dbName, collection string, data ...interface{}) (err error) {
	sc := d.session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Insert(data...)
	return
}

func (d *Database) GetByID(dbName, collection string, id bson.ObjectId, response interface{}) (err error) {
	sc := d.session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).FindId(id).One(response)
	return
}

func (d *Database) Get(dbName, collection string, query interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Skip(offset).Limit(limit).All(response)
	return
}
func (d *Database) GetWithSort(dbName, collection string, query interface{}, sort []string, offset, limit int, response interface{}) (err error) {
	sc := d.session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Sort(sort...).Skip(offset).Limit(limit).All(response)
	return
}
func (d *Database) Update(dbName, collection string, query interface{}, update interface{}) (err error) {
	sc := d.session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Update(query, update)
	return
}
func (d *Database) Aggregate(dbName, collection string, query []bson.M, response interface{}) (err error) {
	sc := d.session.Copy()
	defer sc.Close()
	var a []interface{}
	fmt.Println(query)
	err = sc.DB(dbName).C(collection).Pipe(query).All(&a)
	fmt.Println(a)
	return
}
