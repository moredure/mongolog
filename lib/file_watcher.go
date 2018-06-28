package lib

import (
	"github.com/hpcloud/tail"
	"github.com/mikefaraponov/mongolog/models"
	"log"
	"sync"
)

type Parser interface {
	Parse(string, string) (*models.Log, error)
}

type FileWatcher struct {
	LogEntryParser Parser
	Logs           chan *models.Log
	LogFormat      string
}

func (f *FileWatcher) Watch(filename string, exit chan struct{}, wg *sync.WaitGroup) {
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
			logEntry, err := f.LogEntryParser.Parse(line.Text, filename)
			if err != nil {
				log.Println("Smth bad happened with log entry!")
				break
			}
			f.Logs <- logEntry
		case <-exit:
			t.Stop()
			return
		}
	}
}

func NewFileWatcher(logParser Parser, logFormat string) *FileWatcher {
	return &FileWatcher{
		Logs:           make(chan *models.Log),
		LogEntryParser: logParser,
		LogFormat:      logFormat,
	}
}
