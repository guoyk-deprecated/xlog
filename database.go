package xlog

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"errors"
)

var (
	// IndexedFields fields needed to be indexed
	IndexedFields = []string{"timestamp", "hostname", "env", "project", "topic", "crid"}

	// DistinctFields fields can be queried as distinct
	DistinctFields = []string{"hostname", "env", "project", "topic"}

	// ErrInvalidField field is invalid
	ErrInvalidField = errors.New("invalid field")
)

// Database is a wrapper of mgo.Database
type Database struct {
	DB *mgo.Database

	// CollectionPrefix prefix of collections in database
	CollectionPrefix string
}

// DialDatabase connect a mongo database with options
func DialDatabase(opt Options) (db *Database, err error) {
	var session *mgo.Session
	if session, err = mgo.Dial(opt.Mongo.URL); err != nil {
		return
	}
	db = &Database{
		DB:               session.DB(opt.Mongo.DB),
		CollectionPrefix: opt.Mongo.Collection,
	}
	return
}

// Close close the underlying session
func (d *Database) Close() {
	d.DB.Session.Close()
}

// CollectionName get collection name by date
func (d *Database) CollectionName(t time.Time) string {
	return fmt.Sprintf("%s%04d%02d%02d", d.CollectionPrefix, t.Year(), t.Month(), t.Day())
}

// Collection get collection by date
func (d *Database) Collection(t time.Time) *mgo.Collection {
	return d.DB.C(d.CollectionName(t))
}

// Insert insert a RecordConvertible, choose collection automatically
func (d *Database) Insert(rc RecordConvertible) (err error) {
	var r Record
	if r, err = rc.ToRecord(); err != nil {
		return
	}
	return d.Collection(r.Timestamp).Insert(&r)
}

// EnableSharding enable sharding on collection of the day
func (d *Database) EnableSharding(t time.Time) error {
	return d.DB.Run(bson.D{
		{
			Name:  "shardCollection",
			Value: d.DB.Name + "." + d.CollectionName(t),
		},
		{
			Name:  "key",
			Value: bson.M{"timestamp": "hashed"},
		},
	}, nil)
}

// EnsureIndexes ensure indexes for collection of the day
func (d *Database) EnsureIndexes(t time.Time) (err error) {
	coll := d.Collection(t)
	for _, field := range IndexedFields {
		if err = coll.EnsureIndex(mgo.Index{
			Key:        []string{field},
			Background: true,
		}); err != nil {
			return
		}
	}
	return
}
