package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gennadis/shorturl/internal/config"
	"github.com/gennadis/shorturl/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	ContentType = "Content-Type"
	Location    = "Location"
	PlainText   = "text/plain"
)

type Server struct {
	Store  storage.Repository
	Router *chi.Mux
}

func (s *Server) Start() (err error) {
	s.MountHandlers()
	return http.ListenAndServe(config.Addr, s.Router)
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

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

func (s *Server) expand(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	orignalURL, err := s.Store.Read(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(Location, orignalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func New(storage storage.Repository) *Server {
	return &Server{
		Store:  storage,
		Router: chi.NewRouter(),
	}
}
