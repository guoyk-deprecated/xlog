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

	// Jan 01, 02:00 of the year
	date := time.Date(year, time.January, 1, 2, 0, 0, 0, time.UTC)

	// iterate whole year
	for {
		log.Println("allocating:", date)
		// shard
		if err = db.EnableSharding(date); err != nil {
			log.Println("failed to enable sharding:", date, err)
		}
		// index
		if err = db.EnsureIndexes(date); err != nil {
			log.Println("failed to ensure indexes:", date, err)
		}
		// next day
		date = date.Add(time.Hour * 24)
		// break if not this year
		if date.Year() != year {
			break
		}
	}
}
