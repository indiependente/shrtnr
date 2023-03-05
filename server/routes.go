package server

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	// URLShortenPath is the path used to perform CRUD ops on urls.
	URLShortenPath = `/url`
	// URLResolvePath is the path used to resolve shortened urls.
	URLResolvePath = `/r`
)

func (srv HTTPServer) middlewares() {
	srv.router.Use()
	srv.router.Use(middleware.Compress(5, "application/json"))
	srv.router.Use(middleware.RealIP)
	srv.router.Use(middleware.RequestID)
	srv.router.Use(middleware.Recoverer)
	srv.router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{}))

}

func (srv HTTPServer) routes() {
	srv.router.Handle("/*", http.FileServer(srv.assets))
	srv.router.Get(URLShortenPath+"/:slug", srv.getURL())
	srv.router.Put(URLShortenPath, srv.putURL())
	srv.router.Delete(URLShortenPath+"/:slug", srv.delURL())
	srv.router.Get(URLResolvePath+"/:slug", srv.resolveURL())
	srv.router.Post(URLShortenPath, srv.shortenURL())
}
