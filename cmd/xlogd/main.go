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
		log.Println("failed to load config,", err)
		return
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
	// redis input
	var ri inputs.Input
	if ri, err = inputs.DialRedisInput(options.Redis.URLs[idx], options.Redis.Key); err != nil {
		return
	}
	defer ri.Close()
	// database
	var db *xlog.Database
	if db, err = xlog.DialDatabase(options); err != nil {
		return
	}
	defer db.Close()
	// main loop
	for {
		// check shutdown flag, clear err
		if shutdownMark {
			return
		}
		// read next event
		var rc xlog.RecordConvertible
		if rc, err = ri.Next(); err != nil {
			// redis went wrong, stop input routine
			return
		}
		// skip if timeout
		if rc == nil {
			continue
		}
		// insert document
		if err = db.Insert(rc); err != nil {
			// resend with RPUSH and return, unless it's a conversion error
			if !xlog.IsRecordConversionError(err) {
				ri.Recover(rc)
				// stop on failed
				return
			}
		}
		// increase counter
		atomic.AddUint64(&counter, 1)
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
