package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gennadis/shorturl/internal/app"
)

func (s *Server) shortenPlainText(w http.ResponseWriter, r *http.Request) {
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

	newHash := app.GenerateHash(newURL.String(), s.Config.HashLen)
	if err := s.Store.Write(newHash, newURL.String()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("%s/%s", s.Config.BaseURL, newHash)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}
