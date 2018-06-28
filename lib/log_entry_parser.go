package lib

import (
	"strings"
	"time"
	"github.com/mikefaraponov/mongolog/models"
)

const (
	LOG_TIME = iota
	LOG_MSG  = iota
)

var TimeLayoutByFormat = map[string]string{
	"first_format":  "Feb 1, 2018 at 3:04:05pm (UTC)",
	"second_format": time.RFC3339,
}

type LogEntryParser struct {
	Layout string
	Format string
}

func (t *LogEntryParser) Parse(row, filename string) (*models.Log, error) {
	s := strings.Split(row, " | ")
	logTime, err := time.Parse(t.Layout, s[LOG_TIME])
	if err != nil {
		return nil, err
	}
	return &models.Log{
		logTime,
		s[LOG_MSG],
		filename,
		t.Format,
	}, nil
}

func NewLogEntryParser(layout, format string) *LogEntryParser {
	return &LogEntryParser{layout, format}
}
