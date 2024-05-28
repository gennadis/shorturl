package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gennadis/shorturl/internal/app"
	"github.com/gennadis/shorturl/internal/app/config"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfiguration()
	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("error creating app: %v", err)
	}

	wg := a.BackgroundDeleter.Run(ctx)
	go func() {
		defer close(a.BackgroundDeleter.DeleteChan)
		defer close(a.BackgroundDeleter.ErrorChan)
		wg.Wait()
	}()

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, a.Handler.Router))
}
