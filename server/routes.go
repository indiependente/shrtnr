package server

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/google/uuid"
)

const (
	// URLPath is the path used to shorten urls or resolve url slugs.
	URLPath = `/url`
)

func (srv HTTPServer) middlewares() {
	srv.app.Use(middleware.Compress())
	srv.app.Use(middleware.Recover())
	srv.app.Use(middleware.Pprof())
	srv.app.Use(middleware.RequestID(func() string {
		return uuid.New().String()
	}))
	srv.app.Use(srv.RequestLogger)
}

func (srv HTTPServer) routes() {
	srv.app.Get(URLPath+"/:slug", getURL(srv.svc))
	srv.app.Put(URLPath, putURL(srv.svc))
}

// RequestLogger logs the request
func (srv HTTPServer) RequestLogger(c *fiber.Ctx) {
	method := string(c.Fasthttp.Request.Header.Method())
	url := string(c.Fasthttp.Request.Header.RequestURI())
	srv.log.Event(method).Info(url)
	c.Next()
}
