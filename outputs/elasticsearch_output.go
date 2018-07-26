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

// BulkInsert insert multiple records with ID already set
func (e *ElasticSearch) BulkInsert(rs []xlog.Record) (err error) {
	if len(rs) == 0 {
		return
	}
	bs := e.Client.Bulk()
	for _, r := range rs {
		// clear the ID, let elasticsearch generate it
		r.ID = ""
		// save index name
		name := e.IndexName(r)
		// offset timestamp due to time zone
		r.Timestamp = r.Timestamp.Add(time.Hour * time.Duration(e.TimeOffset))
		// add a bulk operation, REMEMBER do not use pointer for Doc(r)
		bs = bs.Add(elastic.NewBulkIndexRequest().Index(name).Type("_doc").Doc(r))
	}
	_, err = bs.Do(context.Background())
	return
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
