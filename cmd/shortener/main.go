package main

import (
	"log"

	"github.com/gennadis/shorturl/internal/app/config"
	"github.com/gennadis/shorturl/internal/app/server"
	"github.com/gennadis/shorturl/internal/app/storage"
)

func main() {
	cfg := config.NewConfig()
	storage := storage.NewStorage()
	server := server.NewServer(*cfg, storage)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
