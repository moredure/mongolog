package main

import (
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	Session                   *mgo.Session
	Stop                      chan os.Signal
	FileTrackersExitWaitGroup *sync.WaitGroup
	Tracker                   *FileTracker
}

func main() {
	var (
		format             string
		logsCollectionName string
		exitSync           []chan struct{}
		timeLayoutByFormat = map[string]string{
			"first_format":  "Feb 1, 2018 at 3:04:05pm (UTC)",
			"second_format": time.RFC3339,
		}
		config = DefaultConfig()
		stop   = make(chan os.Signal)
		wg     sync.WaitGroup
	)
	defer close(stop)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	flag.StringVar(&format, "format", "first_format", "a first_format or second_format")
	flag.StringVar(&logsCollectionName, "collection", "logs", "a logs mongodb collection name")
	flag.Parse()

	files := flag.Args()

	logParser := NewLogParser(timeLayoutByFormat[format])

	tracker := NewFileTracker(logParser, format)

	for _, filename := range files {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}
		stopTracking := make(chan struct{})
		exitSync = append(exitSync, stopTracking)
		wg.Add(1)
		go tracker.Track(filename, stopTracking, &wg)
	}

	session, err := mgo.Dial(config.MongoUrl)
	if err != nil {
		log.Println("Error establishing mongodb connection!")
		panic(err)
	}
	defer session.Close()

	logs := session.DB("").C(logsCollectionName)

	for {
		select {
		case entry := <-tracker.Logs:
			logs.Insert(&entry)
		case <-stop:
			for _, e := range exitSync {
				close(e)
			}
			return
		}
	}
	wg.Wait()
}
