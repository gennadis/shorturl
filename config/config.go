package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	HashLen    int    `env:"HASH_LEN" envDefault:"7"`
}

func New() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
