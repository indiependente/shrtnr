package server

import (
	"context"

	"github.com/gofiber/fiber"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/shrtnr/service"
)

// HTTPServer implements a Server capable of serving HTTP requests.
type HTTPServer struct {
	app  *fiber.App
	svc  service.Service
	port int
	log  logger.Logger
}

// NewHTTPServer returns a new instance of an HTTPServer.
func NewHTTPServer(app *fiber.App, svc service.Service, port int, log logger.Logger) HTTPServer {
	return HTTPServer{
		app:  app,
		svc:  svc,
		port: port,
		log:  log,
	}
}

// Start starts the HTTP server.
func (srv HTTPServer) Start(ctx context.Context) error {
	return srv.app.Listen(srv.port)
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
