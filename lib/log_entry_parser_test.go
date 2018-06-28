package lib

import (
  "testing"
  "time"
)

var (
  TIME_IN_FIRST_FORMAT_AS_TIME = time.Date(2018, 1, 1, 23, 4, 5, 0, time.UTC)
  TIME_IN_SECOND_FORMAT_AS_TIME = time.Date(2018, 2, 1, 15, 4, 5, 0, time.UTC)
)

const (
  FIRST_FORMAT = "first_format"
  SECOND_FORMAT = "second_format"
  TIME_IN_FIRST_FORMAT = "Feb 1, 2018 at 11:04:05pm (UTC)"
  TIME_IN_SECOND_FORMAT = "2018-02-01T15:04:05Z"
  MSG = "This is log message"
  FILENAME = "/var/log/test.log"
  ROW_IN_FIRST_FORMAT = TIME_IN_FIRST_FORMAT + " | " + MSG
  ROW_IN_SECOND_FORMAT = TIME_IN_SECOND_FORMAT + " | " + MSG
)

func TestFirstFormatParsePass(t *testing.T) {
  firstTimeLayout := TimeLayoutByFormat[FIRST_FORMAT]

  lep := NewLogEntryParser(firstTimeLayout, FIRST_FORMAT)
  log, err := lep.Parse(ROW_IN_FIRST_FORMAT, FILENAME)
  if err != nil {
    t.Error(err)
  }
  if log == nil {
    t.Fail()
  }
  if log.LogTime != TIME_IN_FIRST_FORMAT_AS_TIME {
    t.Error(log.LogTime, TIME_IN_FIRST_FORMAT_AS_TIME)
  }
}

func TestSecondFormatParsePass(t *testing.T) {
  secondTimeLayout := TimeLayoutByFormat[SECOND_FORMAT]

  lep := NewLogEntryParser(secondTimeLayout, SECOND_FORMAT)
  log, err := lep.Parse(ROW_IN_SECOND_FORMAT, FILENAME)
  if err != nil {
    t.Error(err)
  }
  if log == nil {
    t.Fail()
  }
  if log.LogTime != TIME_IN_SECOND_FORMAT_AS_TIME {
    t.Error(log.LogTime, TIME_IN_SECOND_FORMAT_AS_TIME)
  }
}
