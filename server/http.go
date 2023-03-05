package server

import (
	"context"
	"net/http"

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
func NewHTTPServer(r chi.Router, svc service.Service, addr string, assets http.FileSystem, log logger.Logger) (HTTPServer, error) {
	return HTTPServer{
		router: r,
		svc:    svc,
		addr:   addr,
		log:    log,
		assets: assets,
	}, nil
}

// Start starts the HTTP server.
func (srv HTTPServer) Start(ctx context.Context) error {
	return http.ListenAndServe(srv.addr, srv.router)
}

// Shutdown stops the HTTP server.
func (srv HTTPServer) Shutdown(ctx context.Context) error {
	return srv.Shutdown(ctx)
}

// Setup applies all the server configurations enabling startup.
func (srv HTTPServer) Setup(context.Context) error {
	srv.middlewares()
	srv.routes()

	return nil
}
