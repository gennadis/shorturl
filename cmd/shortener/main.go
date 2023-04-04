package main

import (
	"fmt"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/storage/memstore"
)

func main() {
	st := memstore.New()
	app := app.New(st)
	if err := app.Start(); err != nil {
		fmt.Println(err)
	}
}
