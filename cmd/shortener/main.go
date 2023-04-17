package main

import (
	"log"

	"github.com/gennadis/shorturl/config"
	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/storage/memstore"
)

func main() {
	cfg := config.New()
	storage := memstore.New()
	server := app.New(*cfg, storage)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
