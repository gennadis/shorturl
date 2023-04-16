package main

import (
	"fmt"
	"log"

	"github.com/gennadis/shorturl/config"
	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/storage/memstore"
)

func main() {
	cfg := config.New()
	storage := memstore.New()
	server := app.New(*cfg, storage)
	fmt.Print(server.Config.ServerAddr)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
