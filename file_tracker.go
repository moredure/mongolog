package main

import (
	"github.com/hpcloud/tail"
	"log"
	"sync"
)

type FileTracker struct {
	LogParser Parser
	Logs      chan *Log
	LogFormat string
}

func (f *FileTracker) Track(filename string, exit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	t, err := tail.TailFile(filename, tail.Config{Follow: true})
	if err != nil {
		log.Println("Error trying to watch the file", filename)
		return
	}
	defer t.Cleanup()
	for {
		select {
		case line := <-t.Lines:
			logTime, msg, err := f.LogParser.Parse(line.Text)
			if err != nil {
				log.Println("Smth bad happened with log entry!")
				break
			}
			f.Logs <- &Log{
				LogTime:   logTime,
				LogMsg:    msg,
				FilePath:  filename,
				LogFormat: f.LogFormat,
			}
		case <-exit:
			t.Stop()
			return
		}
	}
}

func NewFileTracker(logParser Parser, logFormat string) *FileTracker {
	return &FileTracker{
		Logs:      make(chan *Log),
		LogParser: logParser,
		LogFormat: logFormat,
	}
}
