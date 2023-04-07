package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gennadis/shorturl/internal/config"
	"github.com/gennadis/shorturl/internal/hash"
	"github.com/gennadis/shorturl/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Store  storage.Repository
	Router *chi.Mux
}

func New(storage storage.Repository) *Server {
	return &Server{
		Store:  storage,
		Router: chi.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.MountHandlers()
	return http.ListenAndServe(config.Addr, s.Router)
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.Logger)
	s.Router.Post("/", s.shorten)
	s.Router.Get("/{hash}", s.expand)
}

func (s *Server) shorten(w http.ResponseWriter, r *http.Request) {
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
	s.Store.Write(newHash, newURL.String())
	response := fmt.Sprintf("http://%s/%s", config.Addr, newHash)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

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
