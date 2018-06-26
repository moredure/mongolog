import (
  "gopkg.in/mgo.v2"
  "time"
  "github.com/bitsane/parser/resources"
  "github.com/hpcloud/tail"
)

func main() {
  var (
    format string
    logsCollectionName string
  )
  timeLayoutByFormat := map[string]string{
    "first_format": "Feb 1, 2018 at 3:04:05pm (UTC)",
    "second_format": time.RFC3339,
  }
  config := resources.DefaultConfig()
  stop := make(chan os.Signal)
  signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
  flag.StringVar(&format, "format", "first_format", "a first_format or second_format")
  flag.StringVar(&logsCollectionName, "collection", "logs", "a logs mongodb collection name")
  flag.Parse()
  files := flag.Args()
  logParser := parsers.NewLogParser(timeLayoutByFormat[format])
  tracker := trackers.NewFileTracker(logParser)
  exitSync := make(chan interface{})
  defer close(exitSync)
  for file := range files {
    go tracker.Track(file, exitSync)
  }
  session, err := mgo.Dial(config.MongoUrl)
  if err != nil {
    log.Println("Error establishing mongodb connection!")
    panic(err)
  }
  defer session.Close()
  logs := session.DB().C(logsCollectionName)
  for {
    select {
    case log := <-tracker.LogsChan:
      logs.Insert(&log)
    case stop:
      return
    }
  }
}





