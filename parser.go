package main

type Parser interface {
  Parse() (*time.Time, string, error)
}
