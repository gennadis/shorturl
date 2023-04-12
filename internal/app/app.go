package app

import (
	"net/http"

	"github.com/gennadis/shorturl/config"
	"github.com/gennadis/shorturl/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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

	s.Router.Post("/", s.shortenPlaintext)
	s.Router.Get("/{hash}", s.expand)
}

func New(storage storage.Repository) *Server {
	return &Server{
		Store:  storage,
		Router: chi.NewRouter(),
	}
}
