package main

import (
	"flag"
	"fmt"
	"github.com/mikefaraponov/mongolog/lib"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"syscall"
	"os/signal"
)

const (
	MONGO_URL         = "MONGO_URL"
	DB_FROM_MONGO_URL = ""
)

func main() {
	mongoUrl := os.Getenv(MONGO_URL)
	log.Println("Checking MONGO_URL!")
	if mongoUrl == "" {
		fmt.Println("Setup", MONGO_URL, "environment variable!")
		return
	}
	log.Println("Setting up os exit signal handler!")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	defer close(exit)

	log.Println("Setting up flag parameters!")
	logFormat := flag.String("format", "first_format", "a first_format or second_format")
	logsCollectionName := flag.String("collection", "logs", "a logs mongodb collection name")
	flag.Parse()

	log.Println("Setting up timeLayout!")
	timeLayout, ok := lib.TimeLayoutByFormat[*logFormat]
	if !ok {
		flag.PrintDefaults()
		return
	}

	log.Println("Setting up mongodb connection!")
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Println("Error establishing mongodb connection!", mongoUrl)
		return
	}
	defer session.Close()

	files := flag.Args()

	if len(files) == 0 {
		flag.PrintDefaults()
		return
	}

	log.Println("Setting up mongodb logs collection!")
	logsCollection := session.DB(DB_FROM_MONGO_URL).C(*logsCollectionName)
	log.Println("Setting up log entry parser!")
	logEntryParser := lib.NewLogEntryParser(timeLayout, *logFormat)
	log.Println("Setting up filewatcher!")
	fileWatcher := lib.NewFileWatcher(logEntryParser, *logFormat)
	log.Println("Setting up watchers scheduler!")
	scheduler := &lib.WatchersScheduler{
		LogsCollection: logsCollection,
		FileWatcher:    fileWatcher,
		Exit:           exit,
		Files:          files,
	}
	log.Println("Scheduling watchers!")
	scheduler.Schedule()
}
