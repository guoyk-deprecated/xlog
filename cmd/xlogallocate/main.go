package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/yankeguo/xlog"
)

var (
	options xlog.Options
)

func main() {
	// options
	var err error
	var year int
	flag.IntVar(&year, "year", 0, "allocate indexes for year")
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	if year == 0 {
		panic(errors.New("invalid year input"))
	}

	// create client
	var db *xlog.Database
	if db, err = xlog.DialDatabase(options); err != nil {
		panic(err)
	}

	// 02:00 of the first day of the year
	date := time.Date(year, time.January, 1, 2, 0, 0, 0, time.UTC)
	// iterate whole year
	for {
		coll := db.Collection(date)
		log.Println("allocating:", coll.C.Name)
		// shard
		if err = coll.EnableSharding(); err != nil {
			log.Println("failed to enable sharding:", coll.C.Name, err)
		}
		// index
		if err = coll.EnsureIndexes(); err != nil {
			log.Println("failed to ensure indexes:", coll.C.Name, err)
		}
		// next day
		date = date.Add(time.Hour * 24)
		// break if not this year
		if date.Year() != year {
			break
		}
	}
}
