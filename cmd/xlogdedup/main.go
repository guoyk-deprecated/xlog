/**
 * xlogdedup delete duplicated log entry by checking it's 'timestamp' and 'crid'
 */

package main

import (
	"flag"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/outputs"
)

var (
	date string

	options xlog.Options
)

func main() {
	var err error
	flag.StringVar(&date, "date", "", "date to process, for example '20180720'")
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		log.Println("invalid config,", err)
		return
	}

	var d time.Time
	if d, err = time.Parse("20060102", date); err != nil {
		log.Println("invalid date,", err)
		return
	}

	var db *outputs.MongoDB
	if db, err = outputs.DialMongoDB(options); err != nil {
		log.Println("failed to dial mongodb,", err)
		return
	}
	defer db.Close()

	coll := db.Collection(d)

	// aggregate
	it := coll.Pipe([]bson.M{
		{
			"$group": bson.M{
				"_id":   bson.D{{Name: "timestamp", Value: "$timestamp"}, {Name: "crid", Value: "$crid"}},
				"dups":  bson.M{"$addToSet": "$_id"},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$match": bson.M{"count": bson.M{"$gt": 1}},
		},
	}).AllowDiskUse().Iter()

	// delete duplicates
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
				coll.Remove(bson.M{"_id": id})
			}
		}
	}
}
