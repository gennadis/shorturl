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

type app struct {
	appStorage storage.Repository
}

func New(storage storage.Repository) *app {
	return &app{
		appStorage: storage,
	}
}

func (a *app) Start() error {
	http.HandleFunc("/", a.Multiplex)
	return http.ListenAndServe(config.Addr, nil)
}

func (a *app) Multiplex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.shorten(w, r)
	case http.MethodGet:
		a.expand(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
}

func (a *app) shorten(w http.ResponseWriter, r *http.Request) {
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
	newHash := hash.Generate(config.HashLen)
	a.appStorage.Write(newHash, newURL.String())
	response := fmt.Sprintf("http://%s/%s", config.Addr, newHash)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func (a *app) expand(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[1:] // omit the `/` symbol
	orignalURL, err := a.appStorage.Read(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", orignalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
