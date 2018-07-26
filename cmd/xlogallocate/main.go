package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/outputs"
)

var (
	year    int
	shard   bool
	options xlog.Options
)

func main() {
	// options
	var err error
	flag.IntVar(&year, "year", 0, "allocate indexes for year")
	flag.BoolVar(&shard, "shard", false, "enable sharding")
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	// validate year
	if year == 0 {
		panic(errors.New("invalid year input"))
	}

	// force mongodb tough mode
	options.Mongo.Tough = true

	// create client
	var db *outputs.MongoDB
	if db, err = outputs.DialMongoDB(options); err != nil {
		panic(err)
	}

	// Jan 01, 02:00 of the year
	date := time.Date(year, time.January, 1, 2, 0, 0, 0, time.UTC)

	// iterate whole year
	for {
		log.Println("allocating:", date)
		if shard {
			// shard
			if err = db.EnableSharding(date); err != nil {
				log.Println("failed to enable sharding:", date, err)
				break
			}
		}
		// index
		if err = db.EnsureIndexes(date); err != nil {
			log.Println("failed to ensure indexes:", date, err)
			break
		}
		// next day
		date = date.Add(time.Hour * 24)
		// break if not this year
		if date.Year() != year {
			break
		}
	}
}
