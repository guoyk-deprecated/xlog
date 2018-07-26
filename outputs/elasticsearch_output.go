package outputs

import (
	"github.com/olivere/elastic"
	"fmt"
	"context"
	"github.com/yankeguo/xlog"
	"time"
)

// ElasticSearch wrapper for elastic.Client
type ElasticSearch struct {
	Client *elastic.Client

	TimeOffset int
}

func DialElasticSearch(opt xlog.Options) (e *ElasticSearch, err error) {
	var c *elastic.Client
	if c, err = elastic.NewClient(elastic.SetURL(opt.ES.URLs...)); err != nil {
		return
	}
	e = &ElasticSearch{Client: c, TimeOffset: opt.ES.TimeOffset}
	return
}

// IndexName choose proper index name for record
func (e *ElasticSearch) IndexName(r xlog.Record) string {
	return fmt.Sprintf("%s-%04d-%02d-%02d", r.Topic, r.Timestamp.Year(), r.Timestamp.Month(), r.Timestamp.Day())
}

// Insert insert a record
func (e *ElasticSearch) Insert(rc xlog.RecordConvertible) (err error) {
	var r xlog.Record
	// convert to xlog.Record
	if r, err = rc.ToRecord(); err != nil {
		return
	}
	// save index name
	name := e.IndexName(r)
	// update timestamp due to time zone offset
	r.Timestamp = r.Timestamp.Add(time.Hour * time.Duration(e.TimeOffset))
	// insert document, _doc is the standard document type
	_, err = e.Client.Index().Index(name).Type("_doc").BodyJson(&r).Do(context.Background())
	return
}

// Close
func (e *ElasticSearch) Close() error {
	return nil
}
