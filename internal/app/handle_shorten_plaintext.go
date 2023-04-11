package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gennadis/shorturl/internal/config"
)

func (s *Server) shorten_plaintext(w http.ResponseWriter, r *http.Request) {
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

	newHash := generateHash(config.HashLen)
	if err := s.Store.Write(newHash, newURL.String()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("http://%s/%s", config.Addr, newHash)

	w.Header().Set(ContentType, PlainText)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}
