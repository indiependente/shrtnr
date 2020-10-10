package server

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

const (
	// URLShortenPath is the path used to perform CRUD ops on urls.
	URLShortenPath = `/url`
	// URLResolvePath is the path used to resolve shortened urls.
	URLResolvePath = `/resolve`
)

func (srv HTTPServer) middlewares() {
	srv.app.Use(compress.New())
	srv.app.Use(recover.New())
	srv.app.Use(pprof.New())
	srv.app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// Default middleware config
	srv.app.Use(logger.New(logger.Config{
		Format: "{\"time\": \"${time}\", \"referer\": \"${referer}\", \"protocol\": \"${protocol}\"," +
			" \"ip\": \"${ip}\", \"host\": \"${host}\", \"method\": \"${method}\", \"header\":\"${header:x-request-id}\"," +
			"\"url\": \"${url}\", \"ua\": \"${ua}\", \"latency\": \"${latency}\", \"status\": \"${status}\", \"body\": \"${body}\", " +
			"\"bytesSent\": \"${bytesSent}\", \"bytesReceived\": \"${bytesReceived}\", \"route\": \"${route}\", \"error\": \"${error}\"}\n",
		TimeFormat: time.RFC3339Nano,
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

}

func (srv HTTPServer) routes() {
	srv.app.Get(URLShortenPath+"/:slug", getURL(srv.svc))
	srv.app.Put(URLShortenPath, putURL(srv.svc))
	srv.app.Delete(URLShortenPath+"/:slug", delURL(srv.svc))
	srv.app.Get(URLResolvePath+"/:slug", resolveURL(srv.svc))
	srv.app.Post(URLShortenPath, shortenURL(srv.svc))
}
