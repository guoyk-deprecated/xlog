package xlog

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	// QueryTimeLayouts supported time layouts in begin field
	QueryTimeLayouts = []string{
		"15:04",
		"1504",
		"15:04:05",
		"150405",
	}
)

// Query collection query
type Query struct {
	Begin    time.Time // time start, only hour/minute/second
	Crid     string    // crid
	Hostname string    // hostname
	Env      string    // nev
	Project  string    // project
	Topic    string    // topic
}

// Decode decode a query from request
func (q *Query) Decode(req *http.Request) (err error) {
	if err = req.ParseForm(); err != nil {
		return
	}
	q.Crid = strings.TrimSpace(req.Form.Get("crid"))
	q.Hostname = strings.TrimSpace(req.Form.Get("hostname"))
	q.Env = strings.TrimSpace(req.Form.Get("env"))
	q.Project = strings.TrimSpace(req.Form.Get("project"))
	q.Topic = strings.TrimSpace(req.Form.Get("topic"))

	begin := strings.TrimSpace(req.Form.Get("begin"))
	if len(begin) > 0 {
		for _, layout := range QueryTimeLayouts {
			if q.Begin, err = time.Parse(layout, begin); err == nil {
				break
			}
		}
		err = nil
	}
	return
}

// Execute execute query against mongodb, a date must be specified
func (q *Query) Execute(d *Database, t time.Time, ret *[]LogEntry) error {
	c := d.Collection(t)
	p := bson.M{}
	q.EncodeQuery(p, t)
	sort := "timestamp"
	if q.Begin.IsZero() {
		sort = "-" + sort
	}
	return c.Find(p).Sort(sort).Limit(500).All(ret)
}

// Count count against mongodb, a date must be specified
func (q *Query) Count(d *Database, t time.Time, count *int) (err error) {
	c := d.Collection(t)
	p := bson.M{}
	q.EncodeQuery(p, t)
	*count, err = c.Find(p).Count()
	return
}

// EncodeQuery encode bson.M as mongodb query parameters
func (q Query) EncodeQuery(p bson.M, t time.Time) {
	if len(q.Crid) > 0 {
		p["crid"] = q.Crid
	}
	if len(q.Hostname) > 0 {
		p["hostname"] = q.Hostname
	}
	if len(q.Env) > 0 {
		p["env"] = q.Env
	}
	if len(q.Project) > 0 {
		p["project"] = q.Project
	}
	if len(q.Topic) > 0 {
		p["topic"] = q.Topic
	}
	if !q.Begin.IsZero() {
		p["timestamp"] = bson.M{
			"$gt": time.Date(
				t.Year(),
				t.Month(),
				t.Day(),
				q.Begin.Hour(),
				q.Begin.Minute(),
				q.Begin.Second(),
				0,
				time.UTC,
			),
		}
	}
	return
}

// Encode encode to url query
func (q Query) Encode() string {
	vals := url.Values{}
	if len(q.Crid) > 0 {
		vals.Set("crid", q.Crid)
	}
	if len(q.Hostname) > 0 {
		vals.Set("hostname", q.Hostname)
	}
	if len(q.Env) > 0 {
		vals.Set("env", q.Env)
	}
	if len(q.Project) > 0 {
		vals.Set("project", q.Project)
	}
	if len(q.Topic) > 0 {
		vals.Set("topic", q.Topic)
	}
	if !q.Begin.IsZero() {
		vals.Set("begin", q.BeginFormatted())
	}
	return vals.Encode()
}

// BeginFormatted hh:mm:ss format of begin
func (q Query) BeginFormatted() string {
	if q.Begin.IsZero() {
		return ""
	}
	return fmt.Sprintf("%02d:%02d:%02d", q.Begin.Hour(), q.Begin.Minute(), q.Begin.Second())
}
