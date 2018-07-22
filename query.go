package xlog

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

// TimeRange time range
type TimeRange struct {
	Beginning *time.Time `json:"beginning,omitempty"` // time start
	End       *time.Time `json:"end,omitempty"`       // time end
}

// Query query
type Query struct {
	Timestamp *TimeRange `json:"timestamp,omitempty"` // timestamp
	Crid      string     `json:"crid,omitempty"`      // crid
	Hostname  string     `json:"hostname,omitempty"`  // hostname
	Env       string     `json:"env,omitempty"`       // env
	Project   string     `json:"project,omitempty"`   // project
	Topic     string     `json:"topic,omitempty"`     // topic

	Offset    int  `json:"offset" bson:"-"`    // offset
	Ascendant bool `json:"ascendant" bson:"-"` // timestamp ascendant, default to false
}

// Validated returns a query with Timestamp fixed
func (q Query) Validated() (n Query) {
	n = q
	// assign Timestamp if missing
	if n.Timestamp == nil {
		n.Timestamp = &TimeRange{}
	}
	// assign Beginning as beginning of today if missing
	if n.Timestamp.Beginning == nil {
		b := BeginningOfDay()
		n.Timestamp.Beginning = &b
	}
	// assign End as end of today if missing
	if n.Timestamp.End == nil {
		e := EndOfDay()
		n.Timestamp.End = &e
	}
	// change End to end of the day of Beginning if End is not same date with Beginning or End is before Beginning
	if n.Timestamp.End.Before(*n.Timestamp.Beginning) || !SameDay(*n.Timestamp.Beginning, *n.Timestamp.End) {
		e := EndOfTheDay(*n.Timestamp.Beginning)
		n.Timestamp.End = &e
	}
	// fix offset
	if n.Offset < 0 {
		n.Offset = 0
	}
	// compact fields
	n.Crid = CompactField(n.Crid)
	n.Hostname = CompactField(n.Hostname)
	n.Env = CompactField(n.Env)
	n.Project = CompactField(n.Project)
	n.Topic = CompactField(n.Topic)
	return
}

// Sort field to sort for mongodb
func (q Query) Sort() string {
	if q.Ascendant {
		return "timestamp"
	} else {
		return "-timestamp"
	}
}

// ToMatch convert to bson.M for query match
func (q Query) ToMatch() (m bson.M) {
	m = bson.M{}
	// build $gte and $lt for timestamp
	if q.Timestamp != nil {
		t := bson.M{}
		if q.Timestamp.Beginning != nil {
			t["$gte"] = *q.Timestamp.Beginning
		}
		if q.Timestamp.End != nil {
			t["$lt"] = *q.Timestamp.End
		}
		m["timestamp"] = t
	}
	BSONPutMatchField(m, "crid", q.Crid)
	BSONPutMatchField(m, "hostname", q.Hostname)
	BSONPutMatchField(m, "env", q.Env)
	BSONPutMatchField(m, "project", q.Project)
	BSONPutMatchField(m, "topic", q.Topic)
	return
}
