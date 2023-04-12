package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) expand(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	orignalURL, err := s.Store.Read(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", orignalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
