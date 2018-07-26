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

	var r xlog.Record
	it := coll.Find(bson.M{}).Sort("timestamp").Batch(1000).Prefetch(0.25).Skip(skip).Iter()

	var count = skip
	for {
		if !it.Next(&r) {
			// other error, break
			if err = it.Err(); err != nil {
				log.Println("failed to iterate,", err)
			}
			break
		}
		// save id for debug
		id := r.ID
		// clear ID
		r.ID = ""
		// insert
		if err = eo.Insert(r); err != nil {
			log.Println("failed to insert", id, err)
			break
		}
		// increase count and update percentage
		count++
		if count%1000 == 0 {
			log.Printf("Progress: %10d/%10d", count, total)
		}
	}
	it.Close()
}
