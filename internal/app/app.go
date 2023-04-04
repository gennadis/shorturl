package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gennadis/shorturl/internal/config"
	"github.com/gennadis/shorturl/internal/hash"
	"github.com/gennadis/shorturl/internal/storage"
)

type App struct {
	appStorage storage.Repository
}

func (a *App) Start() error {
	http.HandleFunc("/", a.Multiplex)
	return http.ListenAndServe(config.Addr, nil)
}

func (a *App) Multiplex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.Shorten(w, r)
	case http.MethodGet:
		a.Expand(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
}

func (a *App) Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newURL, err := url.ParseRequestURI(string(b))
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := hash.GenerateHash(config.HashLen)
	a.appStorage.Write(id, newURL.String())

	response := fmt.Sprintf("http://%s/%s", config.Addr, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func (a *App) Expand(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[1:] // omit the `/` symbol
	orignalURL, err := a.appStorage.Read(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", orignalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func New(storage storage.Repository) *App {
	return &App{
		appStorage: storage,
	}
}
