package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/shrtnr/service"
)

// HTTPServer implements a Server capable of serving HTTP requests.
type HTTPServer struct {
	router chi.Router
	svc    service.Service
	addr   string
	log    logger.Logger
	assets http.FileSystem
	http.Server
}

// NewHTTPServer returns a new instance of an HTTPServer.
func NewHTTPServer(r chi.Router, svc service.Service, addr string, assets http.FileSystem, log logger.Logger) (*HTTPServer, error) {
	s := &HTTPServer{
		router: r,
		svc:    svc,
		addr:   addr,
		log:    log,
		assets: assets,
		Server: http.Server{
			Addr:              ":" + addr,
			ReadHeaderTimeout: time.Second,
		},
	}
	s.middlewares(log)
	s.routes()
	s.Server.Handler = s.router

	return s, nil
}

// Start starts the HTTP server.
func (srv *HTTPServer) Start(ctx context.Context) error {
	return srv.ListenAndServe()
}

// Shutdown stops the HTTP server.
func (srv *HTTPServer) Shutdown(ctx context.Context) error {
	return srv.Server.Shutdown(ctx)
}
