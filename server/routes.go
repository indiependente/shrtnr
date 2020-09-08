package server

import (
	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/google/uuid"
)

const (
	// URLShortenPath is the path used to shorten urls.
	URLShortenPath = `/url`
	// URLResolvePath is the path used to resolve shortened urls.
	URLResolvePath = `/resolve`
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
	srv.app.Get(URLShortenPath+"/:slug", getURL(srv.svc))
	srv.app.Put(URLShortenPath, putURL(srv.svc))
	srv.app.Delete(URLShortenPath+"/:slug", delURL(srv.svc))
	srv.app.Get(URLResolvePath+"/:slug", resolveURL(srv.svc))
}

// RequestLogger logs the request
func (srv HTTPServer) RequestLogger(c *fiber.Ctx) {
	method := string(c.Fasthttp.Request.Header.Method())
	url := string(c.Fasthttp.Request.Header.RequestURI())
	status := c.Fasthttp.Response.Header.StatusCode()
	reqID := c.Fasthttp.Response.Header.Peek("X-Request-Id")
	srv.log.Event(url).StatusCode(status).RequestID(strfmt.UUID(reqID)).Info(method)
	c.Next()
}
