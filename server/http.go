package server

import (
	"context"

	"github.com/gofiber/fiber"
)

// HTTPServer implements a Server capable of serving HTTP requests.
type HTTPServer struct {
	app *fiber.App
}

// NewHTTPServer returns a new instance of an HTTPServer.
func NewHTTPServer() HTTPServer {
	app := fiber.New()
	return HTTPServer{
		app: app,
	}
}

// Start starts the HTTP server.
func (srv HTTPServer) Start(ctx context.Context) error {
	return srv.app.Listen(7000)
}

// Shutdown stops the HTTP server.
func (srv HTTPServer) Shutdown(ctx context.Context) error {
	return srv.app.Shutdown()
}

// Routes sets the server's routes.
func (srv HTTPServer) Routes(ctx context.Context) {
	panic("implement me")
}
