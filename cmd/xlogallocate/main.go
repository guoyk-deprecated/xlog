package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"

	"github.com/yankeguo/xlog/types"
)

var (
	options types.Options

	allocateNextYear bool
)

func main() {
	var optionsFile string
	flag.StringVar(&optionsFile, "c", "/etc/xlog.yml", "config file")
	flag.BoolVar(&allocateNextYear, "nextyear", false, "allocate collections for next year")
	flag.Parse()

	// read options file
	var err error
	if err = types.ReadOptionsFile(optionsFile, &options); err != nil {
		panic(err)
	}

	// create client
	var mClient *mgo.Session
	if mClient, err = mgo.Dial(options.Mongo.URL); err != nil {
		panic(err)
	}

	// first day of the year
	year := time.Now().Year()
	if allocateNextYear {
		year++
	}
	date := time.Date(year, time.January, 1, 2, 0, 0, 0, time.UTC)
	// iterate whole year
	for {
		collection := fmt.Sprintf(
			"%s%04d%02d%02d",
			options.Mongo.Collection,
			date.Year(),
			date.Month(),
			date.Day(),
		)
		log.Println("collection:", collection)
		// shard
		if err = mClient.Run(bson.D{
			bson.DocElem{
				Name:  "shardCollection",
				Value: options.Mongo.DB + "." + collection,
			},
			bson.DocElem{
				Name: "key",
				Value: bson.D{bson.DocElem{
					Name:  "timestamp",
					Value: "hashed",
				}},
			},
		}, nil); err != nil {
			log.Println("failed to shard collection:", collection, err)
		}
		// index
		for _, field := range types.IndexedFields {
			if err = mClient.DB(options.Mongo.DB).C(collection).EnsureIndex(mgo.Index{
				Key:        []string{field},
				Background: true,
			}); err != nil {
				log.Println("failed to create index:", collection, field)
			}
		}
		// next day
		date = date.Add(time.Hour * 24)
		// break if not this year
		if date.Year() != year {
			break
		}
	}
}
