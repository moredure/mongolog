package parsers

type LogParser struct {
  Layout string
}

const (
  LOG_TIME = iota
  LOG_MSG = iota
)

func (t *LogParser) Parse(row string) (*time.Time, string, error) {
  s := strings.Split(row, " | ")
  t, err := time.Parse(t.Layout, s[LOG_TIME])
  if err != nil {
    return nil, "", err
  }
  return t, s[LOG_MSG], nil
}

func NewLogParser(layout string) *LogParser {
  return &LogParser{Layout: layout}
}
