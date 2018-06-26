package main

import (
	"github.com/caarlos0/env"
	"log"
	"os"
)

type AppConfig struct {
	MongoUrl string `env:"MONGO_URL,required"`
}

func DefaultConfig() *AppConfig {
	cfg := AppConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("%s", err)
		os.Exit(1)
	}
	return &cfg
}
