package app

import (
	"net/http"

	"github.com/gennadis/shorturl/internal/app/config"
	"github.com/gennadis/shorturl/internal/app/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config config.Config
	Store  storage.Repository
	Router *chi.Mux
}

func (s *Server) Start() (err error) {
	s.MountHandlers()
	return http.ListenAndServe(s.Config.ServerAddr, s.Router)
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Post("/", s.shortenPlainText)
	s.Router.Post("/api/shorten", s.shortenJSON)
	s.Router.Get("/{hash}", s.expand)
}

func NewServer(config config.Config, storage storage.Repository) *Server {
	return &Server{
		Config: config,
		Store:  storage,
		Router: chi.NewRouter(),
	}
}
