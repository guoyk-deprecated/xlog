package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/inputs"
	"github.com/yankeguo/xlog/outputs"
)

var (
	options xlog.Options // options

	shutdownMark  = false             // mark for shutdown
	shutdownGroup = &sync.WaitGroup{} // WaitGroup for shutdown complete

	counter uint64 // number of entries processed
)

func main() {
	// options flag
	var err error
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	// create goroutines for each redis url
	for i := range options.Redis.URLs {
		go inputRoutineGuarded(i)
	}

	// create goroutine for counter reporting
	go reportRoutine()

	// wait for SIGINT or SIGTERM
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	shutdownMark = true

	// wait for all goroutines complete
	log.Println("exiting...")
	shutdownGroup.Wait()
}

func inputRoutineGuarded(idx int) {
	var err error
	for {
		if err = inputRoutine(idx); err != nil {
			log.Println("input routine failed :", idx, ":", err)
		}
		if shutdownMark {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

func inputRoutine(idx int) (err error) {
	log.Println("input routine created:", idx)
	defer log.Println("input routine exited:", idx)
	// maintain shutdown group
	shutdownGroup.Add(1)
	defer shutdownGroup.Done()
	// input
	var xi inputs.Input
	if xi, err = inputs.DialRedis(options.Redis.URLs[idx], options.Redis.Key); err != nil {
		return
	}
	defer xi.Close()
	// output
	var xo outputs.Output
	if options.Mongo.Enabled {
		if xo, err = outputs.DialMongoDB(options); err != nil {
			return
		}
		log.Println("using mongodb")
	} else {
		if xo, err = outputs.DialElasticSearch(options); err != nil {
			return
		}
		log.Println("using elasticsearch")
	}
	defer xo.Close()
	// main loop
	for {
		// check shutdown flag, clear err
		if shutdownMark {
			return
		}
		// read next event
		var rc xlog.RecordConvertible
		if rc, err = xi.Next(); err != nil {
			// redis went wrong, stop input routine
			return
		}
		// skip if timeout
		if rc == nil {
			continue
		}
		// insert document
		if err = xo.Insert(rc); err != nil {
			// recover with RPUSH and return, unless it's a conversion error
			if !xlog.IsRecordConversionError(err) {
				xi.Recover(rc)
				// stop on failed
				return
			}
			// clear conversion error
			err = nil
		} else {
			// increase counter
			atomic.AddUint64(&counter, 1)
		}
	}
}

func reportRoutine() {
	// period
	d := time.Minute * 5
	if options.Verbose {
		d = time.Second * 5
	}
	// ticker
	t := time.NewTicker(d)
	for {
		<-t.C
		log.Println("events emitted:", counter)
	}
}
