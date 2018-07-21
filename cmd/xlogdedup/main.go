/**
 * xlogdedup delete duplicated log entry by checking it's 'timestamp' and 'crid'
 */

package main

import (
	"flag"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/yankeguo/xlog"
)

var (
	date string

	options xlog.Options
)

func main() {
	var err error
	flag.StringVar(&date, "date", "", "date to process, for example '20180720'")
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	var db *xlog.Database
	if db, err = xlog.DialDatabase(options); err != nil {
		panic(db)
	}
	defer db.Close()

	var d time.Time
	if d, err = time.Parse("20060102", date); err != nil {
		log.Println("invalid date,", err)
		return
	}

	coll := db.Collection(d)
	var it *mgo.Iter
	it = coll.C.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": bson.D{
					bson.DocElem{
						Name:  "timestamp",
						Value: "$timestamp",
					},
					bson.DocElem{
						Name:  "crid",
						Value: "$crid",
					},
				},
				"dups": bson.M{
					"$addToSet": "$_id",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$match": bson.M{
				"count": bson.M{
					"$gt": 1,
				},
			},
		},
	}).AllowDiskUse().Iter()
	for {
		doc := struct {
			Dups []bson.ObjectId `bson:"dups"`
		}{}
		if !it.Next(&doc) {
			if err = it.Close(); err != nil {
				log.Println(it.Err())
			}
			break
		}
		if len(doc.Dups) > 0 {
			// IMPORTANT, skip one, or you will delete all
			ids := doc.Dups[1:]
			// in most circumstances, only one dup, so just range
			for _, id := range ids {
				coll.C.Remove(bson.M{
					"_id": id,
				})
			}
		}
	}
}
