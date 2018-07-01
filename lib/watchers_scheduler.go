package lib

import (
	"os"
	"sync"
	"gopkg.in/mgo.v2"
)

type WatchersScheduler struct {
	Files          []string
	Exit           chan os.Signal
	FileWatcher    *FileWatcher
	LogsCollection *mgo.Collection
	ExitSync       []chan struct{}
	WatchersSync   sync.WaitGroup
}

func (s *WatchersScheduler) Schedule() {
	s.startTailers()
	s.trackLogs()
	s.WatchersSync.Wait()
}

func (s *WatchersScheduler) startTailers() {
	for _, filename := range s.Files {
		s.WatchersSync.Add(1)
		stopTracking := make(chan struct{})
		s.ExitSync = append(s.ExitSync, stopTracking)
		go s.FileWatcher.Watch(filename, stopTracking, &s.WatchersSync)
	}
}

func (s *WatchersScheduler) trackLogs() {
	for {
		select {
		case entry := <-s.FileWatcher.Logs:
			s.LogsCollection.Insert(&entry)
		case <-s.Exit:
			for _, e := range s.ExitSync {
				close(e)
			}
			return
		}
	}
}
