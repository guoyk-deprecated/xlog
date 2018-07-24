package xlog

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// QueryLimit limit of result returned
const QueryLimit = 100

var (
	// PlainIndexedFields fields needed to be indexed
	PlainIndexedFields = []string{"timestamp", "hostname", "env", "project", "topic", "crid"}

	// TextIndexedFields fiends needed to be indexed as text
	TextIndexedFields = []string{"message"}
)

// Database is a wrapper of mgo.Database
type Database struct {
	DB *mgo.Database

	// CollectionPrefix prefix of collections in database
	CollectionPrefix string
}

// DialDatabase connect a mongo database with options
func DialDatabase(opt Options) (db *Database, err error) {
	// parse url
	var di *mgo.DialInfo
	if di, err = mgo.ParseURL(opt.Mongo.URL); err != nil {
		return
	}
	// set read timeout to 0 for tough mode
	if opt.Mongo.Tough {
		di.ReadTimeout = 0
	}
	// create session
	var session *mgo.Session
	if session, err = mgo.DialWithInfo(di); err != nil {
		return
	}
	// wrap
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

// Search execute a query
func (d *Database) Search(q Query) (ret Result, err error) {
	coll := d.Collection(q.Timestamp.Beginning)
	// find
	var records []Record
	if err = coll.Find(q.ToMatch()).Sort(q.Sort()).Skip(q.Skip).Limit(QueryLimit).SetMaxTime(time.Second * 45).All(&records); err != nil {
		return
	}
	if records == nil {
		records = []Record{}
	}
	// build result
	ret.Records = records
	ret.Limit = QueryLimit
	return
}

// Trends calculate trends from a query
func (d *Database) Trends(q Query) (rs []Trend, err error) {
	// find collection
	coll := d.Collection(q.Timestamp.Beginning)
	// execute pipeline
	rs = make([]Trend, 0)
	if err = coll.Pipe([]bson.M{{"$match": q.ToMatch()}, {"$group": q.ToGroup()}}).SetMaxTime(time.Second * 45).All(&rs); err != nil {
		return
	}
	return
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
	// plain indexes
	for _, field := range PlainIndexedFields {
		if err = coll.EnsureIndex(mgo.Index{
			Key:        []string{field},
			Background: true,
		}); err != nil {
			return
		}
	}
	// text indexes, should be compounded, because mongo support only one text index
	keys := make([]string, 0)
	for _, field := range TextIndexedFields {
		keys = append(keys, "$text:"+field)
	}
	if err = coll.EnsureIndex(mgo.Index{
		Key:        keys,
		Background: true,
	}); err != nil {
		return
	}
	return
}
