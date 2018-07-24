package xlog

import (
	"github.com/globalsign/mgo/bson"
)

// Query query
type Query struct {
	Timestamp Period   `json:"timestamp,omitempty"` // timestamp
	Crid      []string `json:"crid,omitempty"`      // crid, using $in
	Hostname  []string `json:"hostname,omitempty"`  // hostname, using $in
	Env       []string `json:"env,omitempty"`       // env, using $in
	Project   []string `json:"project,omitempty"`   // project, using $in
	Topic     []string `json:"topic,omitempty"`     // topic, using $in
	Message   []string `json:"message,omitempty"`   // message content query, basically $and mode

	Skip int `json:"skip" bson:"-"` // skip
}

// Validated returns a query with Timestamp fixed
func (q Query) Validated() (n Query) {
	// assign n with q
	n = q
	// assign Beginning as beginning of today if missing
	if n.Timestamp.Beginning.IsZero() {
		n.Timestamp.Beginning = BeginningOfDay()
	}
	// assign End as end of today if missing
	if n.Timestamp.End.IsZero() {
		n.Timestamp.End = EndOfDay()
	}
	// change End to end of the day of Beginning if End is not same date with Beginning or End is before Beginning
	if n.Timestamp.End.Before(n.Timestamp.Beginning) || !SameDay(n.Timestamp.Beginning, n.Timestamp.End) {
		n.Timestamp.End = EndOfTheDay(n.Timestamp.Beginning)
	}
	// fix offset
	if n.Skip < 0 {
		n.Skip = 0
	}
	return
}

func (q Query) TrendQueries() []Query {
	qs := make([]Query, 0)
	ts := q.Timestamp.TrendPeriods()
	for _, t := range ts {
		nq := q
		nq.Timestamp = t
		qs = append(qs, nq)
	}
	return qs
}

// Sort field to sort for mongodb
func (q Query) Sort() string {
	if q.Timestamp.Ascendant {
		return "timestamp"
	} else {
		return "-timestamp"
	}
}

// ToMatch convert to bson.M for query match
func (q Query) ToMatch() (m bson.M) {
	m = bson.M{}
	m["timestamp"] = q.Timestamp.ToMatch()
	BSONPutMatchField(m, "crid", q.Crid)
	BSONPutMatchField(m, "hostname", q.Hostname)
	BSONPutMatchField(m, "env", q.Env)
	BSONPutMatchField(m, "project", q.Project)
	BSONPutMatchField(m, "topic", q.Topic)
	BSONPutTextField(m, q.Message)
	return
}
