package main

import (
	"fmt"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/storage/memstore"
)

func main() {
	storage := memstore.New()
	shortener := app.New(storage)
	if err := shortener.Start(); err != nil {
		fmt.Println(err)
	}
}
