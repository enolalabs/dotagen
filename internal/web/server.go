package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/enolalabs/dotagen/v2/internal/platform"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	port     int
	rootDir  string
	registry *platform.Registry
	mu       chan struct{}
}

func NewServer(rootDir string, port int) (*Server, error) {
	return &Server{
		port:     port,
		rootDir:  rootDir,
		registry: platform.NewRegistry(),
		mu:       make(chan struct{}, 1),
	}, nil
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to load static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	mux.HandleFunc("GET /api/config", s.handleGetConfig)
	mux.HandleFunc("PUT /api/config", s.handleUpdateConfig)
	mux.HandleFunc("GET /api/agents", s.handleListAgents)
	mux.HandleFunc("GET /api/agents/{name}", s.handleGetAgent)
	mux.HandleFunc("POST /api/agents", s.handleCreateAgent)
	mux.HandleFunc("PUT /api/agents/{name}", s.handleUpdateAgent)
	mux.HandleFunc("DELETE /api/agents/{name}", s.handleDeleteAgent)
	mux.HandleFunc("GET /api/targets", s.handleListTargets)
	mux.HandleFunc("GET /api/preview/{agent}/{target}", s.handlePreview)
	mux.HandleFunc("POST /api/sync", s.handleSync)
	mux.HandleFunc("POST /api/sync/{target}", s.handleSyncTarget)
	mux.HandleFunc("POST /api/clean", s.handleClean)
	mux.HandleFunc("GET /api/status", s.handleStatus)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}

func (s *Server) lock() {
	s.mu <- struct{}{}
}

func (s *Server) unlock() {
	<-s.mu
}
