package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/shrtnr/service"
)

// HTTPServer implements a Server capable of serving HTTP requests.
type HTTPServer struct {
	app    *fiber.App
	svc    service.Service
	port   int
	log    logger.Logger
	assets http.FileSystem
}

// NewHTTPServer returns a new instance of an HTTPServer.
func NewHTTPServer(app *fiber.App, svc service.Service, port int, assets http.FileSystem, log logger.Logger) (HTTPServer, error) {
	return HTTPServer{
		app:    app,
		svc:    svc,
		port:   port,
		log:    log,
		assets: assets,
	}, nil
}

// Start starts the HTTP server.
func (srv HTTPServer) Start(ctx context.Context) error {
	return srv.app.Listen(fmt.Sprintf(":%d", srv.port))
}

// Shutdown stops the HTTP server.
func (srv HTTPServer) Shutdown(ctx context.Context) error {
	return srv.app.Shutdown()
}

// Setup applies all the server configurations enabling startup.
func (srv HTTPServer) Setup(context.Context) error {
	srv.middlewares()
	srv.routes()
	return nil
}
