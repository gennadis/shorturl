package main

import (
	"log"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/app/config"
	"github.com/gennadis/shorturl/internal/app/storage"
)

func main() {
	cfg := config.New()
	log.Printf("filepath env value: %q", cfg.Filepath)
	storage := storage.New()
	server := app.New(*cfg, storage)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
