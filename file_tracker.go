
type FileTracker struct {
  LogParser *Parser
  Logs chan Log
}

func (f *FileTracker) Track(filename string, exit chan interface{}) {
  if _, err := os.Stat(filename); os.IsNotExist(err) {
    return
  }
  t, err := tail.TailFile(filename, tail.Config{Follow: true})
  if err != nil {
    log.Println("Error trying to watch the file", filename)
    panic(err)
  }
  defer t.Cleanup()
  for {
    select {
    case line := <-t.Lines:
      f.Logs <- f.LogParser.Parse(line)
    case <-exit:
      t.Stop()
      return
    }
  }
}

func NewFileTracker(logParser *Parser) {
  return &FileTracker{
    Logs: make(chan *Logs),
    LogParser: logParser,
  }
}
