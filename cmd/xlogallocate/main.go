package main

import (
	"flag"
	"log"
	"time"

	"github.com/yankeguo/xlog"
)

var (
	options xlog.Options

	next bool
)

func main() {
	// options
	var err error
	flag.BoolVar(&next, "next", false, "allocate for next year")
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	// create client
	var db *xlog.Database
	if db, err = xlog.DialDatabase(options); err != nil {
		panic(err)
	}

	// get year
	year := time.Now().Year()
	if next {
		year++
	}
	// 02:00 of the first day of the year
	ts := time.Date(year, time.January, 1, 2, 0, 0, 0, time.UTC)
	// iterate whole year
	for {
		coll := db.CollectionName(ts)
		log.Println("allocating:", coll)
		// shard
		if err = db.EnableSharding(coll); err != nil {
			log.Println("failed to enable sharding:", coll, err)
		}
		// index
		if err = db.EnsureIndexes(coll); err != nil {
			log.Println("failed to ensure indexes:", coll, err)
		}
		// next day
		ts = ts.Add(time.Hour * 24)
		// break if not this year
		if ts.Year() != year {
			break
		}
	}
}
