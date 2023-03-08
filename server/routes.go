package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/indiependente/pkg/logger"
)

const (
	// PingPath is used to check the service is still up and running by load balancers etc.
	PingPath = `/ping`
	// URLShortenPath is the path used to perform CRUD ops on urls.
	URLShortenPath = `/url`
	// URLResolvePath is the path used to resolve shortened urls.
	URLResolvePath = `/r`
	// ProfileDebugPath is the path used to read profiling data.
	ProfileDebugPath        = `/debug`
	defaultCompressionLevel = 5
)

func (srv *HTTPServer) middlewares(log logger.Logger) {
	srv.router.Use(middleware.Heartbeat(PingPath))
	srv.router.Use(middleware.Timeout(time.Minute))
	srv.router.Use(middleware.RealIP)
	srv.router.Use(middleware.RequestID)
	srv.router.Use(FastLogger(log))
	srv.router.Use(middleware.Recoverer)
	srv.router.Use(middleware.URLFormat)
	srv.router.Use(render.SetContentType(render.ContentTypeJSON))
	srv.router.Use(middleware.Compress(defaultCompressionLevel, "application/json"))
	srv.router.Mount(ProfileDebugPath, middleware.Profiler())
}

func (srv *HTTPServer) routes() {
	srv.router.Handle("/*", http.FileServer(srv.assets))
	srv.router.Route(URLShortenPath, func(r chi.Router) {
		r.Put("/", srv.putURL())
		r.Post("/", srv.shortenURL())
		r.Route("/{slug}", func(r chi.Router) {
			r.Get("/", srv.getURL())
			r.Delete("/", srv.delURL())
		})
	})

	srv.router.Get(URLResolvePath+"/{slug}", srv.resolveURL())
}
