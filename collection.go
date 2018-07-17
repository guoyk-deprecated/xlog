package xlog

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	sizeSuffixes = []struct {
		S string
		N int
	}{
		{"TB", 1024 * 1024 * 1024 * 1024},
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"KB", 1024},
		{"B", 1},
	}
)

// CollectionStats stats of database
type CollectionStats struct {
	StorageSize int `bson:"storageSize"`
}

// StorageSizeFormatted formatted storage size
func (c CollectionStats) StorageSizeFormatted() string {
	if c.StorageSize <= 0 {
		return "0"
	}
	for _, s := range sizeSuffixes {
		if c.StorageSize > s.N {
			return fmt.Sprintf("%.02f%s", float64(c.StorageSize)/float64(s.N), s.S)
		}
	}
	return ""
}

// Collection wrapper of mgo.Collection
type Collection struct {
	C      *mgo.Collection
	Date   time.Time // date of the collection, only year/month/day available
	Prefix string    // prefix of collection name
}

type distinctResult struct {
	Values []string `bson:"values"`
}

// Distinct distinct field values
func (c *Collection) Distinct(field string, out *[]string) (err error) {
	var ret distinctResult
	err = c.C.Database.Run(bson.D{
		bson.DocElem{
			Name:  "distinct",
			Value: c.C.Name,
		},
		bson.DocElem{
			Name:  "key",
			Value: field,
		},
	}, &ret)
	*out = ret.Values
	return
}

// Stats get stats of collection
func (c *Collection) Stats(out *CollectionStats) (err error) {
	c.C.Database.Run(bson.D{
		bson.DocElem{
			Name:  "collStats",
			Value: c.C.Name,
		},
	}, out)
	return
}

// Execute execute a query on collection
func (c *Collection) Execute(q Query, out interface{}) (err error) {
	p := bson.M{}
	q.ToBSON(p, c.Date)
	sort := "timestamp"
	if q.Begin.IsZero() || q.End.IsZero() {
		sort = "-" + sort
	}
	return c.C.Find(p).Sort(sort).Limit(200).All(out)
}

// Count count number of logs by query
func (c *Collection) Count(q Query, out *int) (err error) {
	p := bson.M{}
	q.ToBSON(p, c.Date)
	*out, err = c.C.Find(p).Count()
	return
}

// TotalCount total count
func (c *Collection) TotalCount(out *int) (err error) {
	*out, _ = c.C.Count()
	return
}

// EnableSharding enable sharding
func (c *Collection) EnableSharding() error {
	return c.C.Database.Run(bson.D{
		bson.DocElem{
			Name:  "shardCollection",
			Value: c.C.Database.Name + "." + c.C.Name,
		},
		bson.DocElem{
			Name: "key",
			Value: bson.D{bson.DocElem{
				Name:  "timestamp",
				Value: "hashed",
			}},
		},
	}, nil)
}

// EnsureIndexes ensure indexes for collection
func (c *Collection) EnsureIndexes() (err error) {
	for _, field := range IndexedFields {
		if err = c.C.EnsureIndex(mgo.Index{
			Key:        []string{field},
			Background: true,
		}); err != nil {
			return
		}
	}
	return
}
