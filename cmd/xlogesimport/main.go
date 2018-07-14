package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/globalsign/mgo"
	"github.com/olivere/elastic"
	"github.com/yankeguo/xlog"
)

var (
	options xlog.Options

	esIndex    string
	esEndpoint string
	debugMode  bool
)

func main() {
	ctx := context.Background()

	var optionsFile string
	flag.StringVar(&optionsFile, "c", "/etc/xlog.yml", "config file")
	flag.StringVar(&esIndex, "esindex", "", "elasticsearch index to import")
	flag.StringVar(&esEndpoint, "esendpoint", "http://127.0.0.1:9200", "elasticsearch endpoint")
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")
	flag.Parse()

	if len(esIndex) == 0 {
		panic(errors.New("elasticsearch index not specified"))
	}
	if len(esEndpoint) == 0 {
		panic(errors.New("elasticsearch endpoint not specified"))
	}

	var err error

	// read options file
	if err = xlog.ReadOptionsFile(optionsFile, &options); err != nil {
		panic(err)
	}

	// create elasticsearch client
	var esClient *elastic.Client
	if esClient, err = elastic.NewClient(elastic.SetURL(esEndpoint)); err != nil {
		panic(err)
	}

	// create mongo client
	var mClient *mgo.Session
	if mClient, err = mgo.Dial(options.Mongo.URL); err != nil {
		panic(err)
	}
	mDB := mClient.DB(options.Mongo.DB)

	var exists bool
	if exists, err = esClient.IndexExists(esIndex).Do(ctx); err != nil {
		panic(err)
	}
	if !exists {
		panic(errors.New("index " + esIndex + " not exists"))
	}

	var total int64
	if total, err = esClient.Count(esIndex).Do(ctx); err != nil {
		panic(err)
	}

	scroll := esClient.Scroll(esIndex)

	var result *elastic.SearchResult
	var ee xlog.ESEntry
	var le xlog.LogEntry
	var count int64
	for {
		result, err = scroll.Do(ctx)
		if err == nil {
			// ok
		} else if err == io.EOF {
			break
		} else {
			panic(err)
		}
		for _, hit := range result.Hits.Hits {
			if err = json.Unmarshal(*hit.Source, &ee); err != nil {
				log.Println("failed to unmarshal")
				continue
			}
			if err = ee.Convert(&le); err != nil {
				log.Println("failed to convert")
				continue
			}
			collection := fmt.Sprintf(
				"%s%04d%02d%02d",
				options.Mongo.Collection,
				le.Timestamp.Year(),
				le.Timestamp.Month(),
				le.Timestamp.Day(),
			)
			count++
			log.Printf("progress: %012d / %012d", count, total)
			if debugMode {
				log.Println(">>> WILL INSERT", collection)
				for k, v := range le.ToBSON() {
					log.Println("  ", k, "=", v)
				}
			} else {
				if err = mDB.C(collection).Insert(le.ToBSON()); err != nil {
					log.Println(err)
				}
			}
		}
		if debugMode && count > 100 {
			break
		}
	}
}
