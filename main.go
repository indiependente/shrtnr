package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/pkg/shutdown"
	"github.com/indiependente/shrtnr/repository"
	"github.com/indiependente/shrtnr/server"
	"github.com/indiependente/shrtnr/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	appName = "shrtnr"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal("unexpected failure: ", err)
	}
}

func run() error {
	log := logger.GetLoggerString(appName, "DEBUG")

	mongoConf := repository.BuildMongoConfigs()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConf.URI()))
	if err != nil {
		return fmt.Errorf("could not start mongodb: %w", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to mongodb: %w", err)
	}
	defer client.Disconnect(ctx) // nolint: errcheck

	db := client.Database(mongoConf.DB)

	// create store
	coll := db.Collection(mongoConf.Collection)
	store := repository.NewMongoDBURLStorer(coll)

	// create slugger
	slugLen, err := strconv.Atoi(os.Getenv("SLUG_LEN"))
	if err != nil {
		return fmt.Errorf("could not parse SLUG_LEN: %w", err)
	}
	slugger := service.NewFixedLenSlugger(slugLen)
	// create service
	svc := service.NewURLService(store, slugger)
	// create server
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return fmt.Errorf("could not parse PORT: %w", err)
	}
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})
	srv := server.NewHTTPServer(app, svc, port, log)
	err = srv.Setup(ctx)
	if err != nil {
		return fmt.Errorf("error while running server setup: %w", err)
	}

	// Start HTTP server
	go func() {
		err := srv.Start(ctx)
		if err != nil {
			fmt.Printf("%+v\n", err)
			log.Fatal("error while running HTTP server", err)
		}
	}()

	// Wait
	err = shutdown.Wait(ctx, cancel, srv.Shutdown)
	if err != nil {
		return fmt.Errorf("error while shutting down server: %w", err)
	}
	return nil
}
