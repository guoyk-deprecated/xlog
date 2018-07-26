package main

import (
	"github.com/yankeguo/xlog"
	"flag"
	"time"
	"github.com/yankeguo/xlog/outputs"
	"log"
	"github.com/globalsign/mgo/bson"
)

var (
	date string
	skip int

	options xlog.Options
)

// Progress is the progress counter
type Progress struct {
	Total int
	Count int
}

// Increase the counter
func (p *Progress) Increase() {
	p.Count++
	if p.Count%1000 == 0 {
		log.Printf("Progress: %10d/%10d", p.Count, p.Total)
	}
}

func main() {
	var err error
	flag.StringVar(&date, "date", "", "date to process, for example '20180720'")
	flag.IntVar(&skip, "skip", 0, "skip records")
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		log.Println("invalid config,", err)
		return
	}

	var d time.Time
	if d, err = time.Parse("20060102", date); err != nil {
		log.Println("invalid date,", err)
		return
	}

	// force tough mode
	options.Mongo.Tough = true

	var mo *outputs.MongoDB
	if mo, err = outputs.DialMongoDB(options); err != nil {
		log.Println("failed to dial mongodb,", err)
		return
	}
	defer mo.Close()

	coll := mo.Collection(d)

	var total int
	if total, err = coll.Count(); err != nil {
		log.Println("failed to count,", err)
	}

	if total == 0 {
		return
	}

	var eo *outputs.ElasticSearch
	if eo, err = outputs.DialElasticSearch(options); err != nil {
		log.Println("failed to dial elasticsearch,", err)
	}

	p := &Progress{Total: total, Count: skip}

	it := coll.Find(bson.M{}).Sort("timestamp").Batch(1000).Prefetch(0.25).Skip(skip).Iter()
	defer it.Close()

	r := xlog.Record{}
	rs := make([]xlog.Record, 0)

	// main loop
	for {
		// clear the variable
		r = xlog.Record{}
		// request for the next
		if !it.Next(&r) {
			// other error, return
			if err = it.Err(); err != nil {
				log.Println("failed to iterate,", err)
				return
			}
			// break for final Next()
			break
		}
		// increase progress
		p.Increase()
		// append to rs
		rs = append(rs, r)
		// if rs >= 1000, do the bulk insertion
		if len(rs) >= 1000 {
			// do bulk insertion
			if err = eo.BulkInsert(rs); err != nil {
				log.Println("failed to insert", err)
				return
			}
			// clear and reuse the rs
			rs = rs[:0]
		}
	}
	// insert the rest
	if err = eo.BulkInsert(rs); err != nil {
		log.Println("failed to insert", err)
		return
	}
}
