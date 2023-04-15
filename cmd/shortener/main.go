package main

import (
	"log"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/storage/memstore"
)

func main() {
	storage := memstore.New()
	server := app.New(storage)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
