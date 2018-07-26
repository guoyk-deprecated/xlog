package outputs

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/yankeguo/xlog"
)

// QueryLimit limit of result returned
const QueryLimit = 100

var (
	// PlainIndexedFields fields needed to be indexed
	PlainIndexedFields = []string{"timestamp", "hostname", "env", "project", "topic", "crid"}

	// TextIndexedFields fiends needed to be indexed as text
	TextIndexedFields = []string{"message"}
)

// MongoDB is a wrapper of mgo.MongoDB
type MongoDB struct {
	DB *mgo.Database

	// CollectionPrefix prefix of collections in database
	CollectionPrefix string
}

// DialMongoDB connect a mongo database with options
func DialMongoDB(opt xlog.Options) (db *MongoDB, err error) {
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
	if opt.Mongo.Tough {
		session.SetCursorTimeout(0)
	}
	// wrap
	db = &MongoDB{
		DB:               session.DB(opt.Mongo.DB),
		CollectionPrefix: opt.Mongo.Collection,
	}
	return
}

// Close close the underlying session
func (d *MongoDB) Close() error {
	d.DB.Session.Close()
	return nil
}

// CollectionName get collection name by date
func (d *MongoDB) CollectionName(t time.Time) string {
	return fmt.Sprintf("%s%04d%02d%02d", d.CollectionPrefix, t.Year(), t.Month(), t.Day())
}

// Collection get collection by date
func (d *MongoDB) Collection(t time.Time) *mgo.Collection {
	return d.DB.C(d.CollectionName(t))
}

// Insert insert a RecordConvertible, choose collection automatically
func (d *MongoDB) Insert(rc xlog.RecordConvertible) (err error) {
	var r xlog.Record
	if r, err = rc.ToRecord(); err != nil {
		return
	}
	return d.Collection(r.Timestamp).Insert(&r)
}

// Search execute a query
func (d *MongoDB) Search(q xlog.Query) (ret xlog.Result, err error) {
	coll := d.Collection(q.Timestamp.Beginning)
	// find
	var records []xlog.Record
	if err = coll.Find(q.ToMatch()).Sort(q.Sort()).Skip(q.Skip).Limit(QueryLimit).SetMaxTime(time.Second * 45).All(&records); err != nil {
		return
	}
	if records == nil {
		records = []xlog.Record{}
	}
	// build result
	ret.Records = records
	ret.Limit = QueryLimit
	return
}

// Trends calculate trends from a query
func (d *MongoDB) Trends(q xlog.Query) (rs []xlog.Trend, err error) {
	coll := d.Collection(q.Timestamp.Beginning)
	// trend queries
	qs := q.TrendQueries()
	rs = make([]xlog.Trend, 0, len(qs))
	for _, tq := range qs {
		var c int
		if c, err = coll.Find(tq.ToMatch()).SetMaxTime(time.Second * 10).Count(); err != nil {
			return
		}
		rs = append(rs, xlog.Trend{
			Beginning: tq.Timestamp.Beginning,
			End:       tq.Timestamp.End,
			Count:     c,
		})
	}
	return
}

// EnableSharding enable sharding on collection of the day
func (d *MongoDB) EnableSharding(t time.Time) error {
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
func (d *MongoDB) EnsureIndexes(t time.Time) (err error) {
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
