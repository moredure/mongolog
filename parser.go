package main

import (
	"time"
)

type Parser interface {
	Parse(string) (*time.Time, string, error)
}
