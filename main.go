package main

import (
	"context"

	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/pkg/shutdown"
	"github.com/indiependente/shrtnr/server"
)

const appName = "shrtnr"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := logger.GetLoggerString(appName, "DEBUG")

	srv := server.NewHTTPServer()
	// Start HTTP server
	go func() {
		err := srv.Start(ctx)
		if err != nil {
			log.Fatal("error while running HTTP server", err)
		}
	}()

	// Wait
	err := shutdown.Wait(ctx, cancel, srv.Shutdown)
	if err != nil {
		log.Fatal("error while shutting down server", err)
	}
}
