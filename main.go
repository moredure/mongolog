package main

import (
	"flag"
	"fmt"
	"github.com/mikefaraponov/mongolog/lib"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	MONGO_URL         = "MONGO_URL"
	DB_FROM_MONGO_URL = ""
)

func main() {
	mongoUrl := os.Getenv(MONGO_URL)

	if mongoUrl == "" {
		fmt.Println("Setup", MONGO_URL, "environment variable!")
		return
	}

	exit := make(chan os.Signal)
	defer close(exit)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	logFormat := flag.String("format", "first_format", "a first_format or second_format")
	logsCollectionName := flag.String("collection", "logs", "a logs mongodb collection name")
	flag.Parse()

	timeLayout, ok := lib.TimeLayoutByFormat[*logFormat]
	if !ok {
		flag.PrintDefaults()
		return
	}

	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Println("Error establishing mongodb connection!")
		return
	}
	defer session.Close()

	logsCollection := session.DB(DB_FROM_MONGO_URL).C(*logsCollectionName)
	logEntryParser := lib.NewLogEntryParser(timeLayout, *logFormat)
	fileWatcher := lib.NewFileWatcher(logEntryParser, *logFormat)
	scheduler := &lib.WatchersScheduler{
		LogsCollection: logsCollection,
		FileWatcher:    fileWatcher,
		Exit:           exit,
		Files:          flag.Args(),
	}

	scheduler.Schedule()
}
