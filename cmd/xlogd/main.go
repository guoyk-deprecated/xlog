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
	for _, url := range options.Redis.URLs {
		go beaterRoutineLoop(url)
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

func beaterRoutineLoop(redisURL string) {
	var err error
	for {
		if err = beaterRoutine(redisURL); err != nil {
			log.Println("subroutine failed :", redisURL, ":", err)
		}
		if shutdownMark {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

func beaterRoutine(redisURL string) (err error) {
	log.Println("subroutine created:", redisURL)
	defer log.Println("subroutine exited:", redisURL)
	// maintain shutdown group
	shutdownGroup.Add(1)
	defer shutdownGroup.Done()
	// beater
	var bt *xlog.Beater
	if bt, err = xlog.DialBeater(redisURL, options); err != nil {
		return
	}
	defer bt.Close()
	// database
	var db *xlog.Database
	if db, err = xlog.DialDatabase(options); err != nil {
		return
	}
	defer db.Close()
	// variables
	var (
		be xlog.BeatEntry
		le xlog.LogEntry
	)
	// main loop
	for {
		// check shutdown flag
		if shutdownMark {
			return nil
		}
		// read next beat event
		if err = bt.NextEvent(&be); err != nil {
			if err == xlog.ErrBeaterTimeout {
				// timeout is normal
				continue
			} else if err == xlog.ErrBeaterMalformed {
				if options.Verbose {
					log.Println("non-JSON beat event detected")
				}
				// ignore malformed
				continue
			} else {
				// redis went wrong, stop subroutine
				return
			}
		}
		// convert beat entry to log entry
		if !be.Convert(&le) {
			if options.Verbose {
				log.Println("failed to convert beat event:")
				log.Println("  beat.hostname =", be.Beat.Hostname)
				log.Println("  source        =", be.Source)
				log.Println("  message       =", be.Message)
			}
			// ignore on failed
			continue
		}
		// insert document
		if err = db.Insert(le); err != nil {
			// resend with RPUSH
			bt.Recover(be)
			// stop on failed
			return
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
