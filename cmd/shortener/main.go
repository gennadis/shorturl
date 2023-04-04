package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gennadis/shorturl/internal/app/config"
	"github.com/gennadis/shorturl/internal/app/hash"
)

var urls = make(map[string]string)

// Accepts POST requests with `url` to shorten in request body.
// Returns 201 and short url in request body if successful.
func shortenURL(w http.ResponseWriter, r *http.Request) {
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
	urls[id] = newURL.String()

	response := fmt.Sprintf("%s/%s", config.Addr, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[1:] // omit the `/` symbol
	orignalURL, ok := urls[hash]
	if !ok {
		http.Error(w, "Wrong hash provided", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", orignalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func multiplexer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		shortenURL(w, r)
	case http.MethodGet:
		redirectURL(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/", multiplexer)
	http.ListenAndServe(config.Addr, nil)
}
