package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
	"github.com/yankeguo/xlog"
)

var (
	options xlog.Options

	shutdownFlag  = false             // flag for shutdown
	shutdownGroup = &sync.WaitGroup{} // WaitGroup for shutdown complete
)

func main() {
	var optionsFile string
	flag.StringVar(&optionsFile, "c", "/etc/xlog.yml", "config file")
	flag.Parse()

	// read options file
	var err error
	if err = xlog.ReadOptionsFile(optionsFile, &options); err != nil {
		panic(err)
	}

	// create goroutines for each redis url
	for _, url := range options.Redis.URLs {
		go drainRedisWithRetry(url)
	}

	// wait for SIGINT or SIGTERM
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	shutdownFlag = true

	// wait for all goroutines complete
	log.Println("exiting...")
	shutdownGroup.Wait()
}

func drainRedisWithRetry(rURL string) {
	var err error
	for {
		if err = drainRedis(rURL); err != nil {
			log.Println("error occured for", rURL, ":", err)
		}
		if shutdownFlag {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

func drainRedis(rURL string) (err error) {
	log.Println("routine created", rURL)
	defer log.Println("routine exited", rURL)
	// maintain shutdown group
	shutdownGroup.Add(1)
	defer shutdownGroup.Done()
	// redis client
	var rClient *redis.Client
	if rClient, err = createRedisClient(rURL); err != nil {
		return
	}
	defer rClient.Close()
	// mongo client
	var mClient *mgo.Session
	if mClient, err = mgo.Dial(options.Mongo.URL); err != nil {
		return
	}
	defer mClient.Close()
	// mongo database instance
	mDB := mClient.DB(options.Mongo.DB)
	// main loop
	for {
		// exit if shutdown flag is set
		if shutdownFlag {
			break
		}
		// variables
		var bEntry xlog.BeatEntry
		var lEntry xlog.LogEntry
		// blpop redis list element, timeout 2 seconds
		var eStrs []string
		if eStrs, err = rClient.BLPop(
			time.Second*3,
			options.Redis.Key,
		).Result(); err != nil && err != redis.Nil {
			err = nil
			return
		}
		// continue if empty list
		if len(eStrs) == 0 {
			continue
		}
		// for list element
		for i, eStr := range eStrs {
			// unmarshal JSON
			if err = json.Unmarshal([]byte(eStr), &bEntry); err != nil {
				err = nil
				continue
			}
			// convert BeatEntry to LogEntry
			if !bEntry.Convert(&lEntry) {
				continue
			}
			// decide collection name
			var cName = lEntry.CollectionName(options.Mongo.Collection)
			if err = mDB.C(cName).Insert(lEntry.ToBSON()); err != nil {
				// RPUSH remaining strings
				rClient.RPush(options.Redis.Key, eStr[i:])
				return
			}
		}
	}
	return
}

// create a redis client with server pinged
func createRedisClient(url string) (rClient *redis.Client, err error) {
	var rOpt *redis.Options
	if rOpt, err = redis.ParseURL(url); err != nil {
		// panic if url is malformed
		panic(err)
	}
	// create redis client and ping
	rClient = redis.NewClient(rOpt)
	if err = rClient.Ping().Err(); err != nil {
		return
	}
	return
}
