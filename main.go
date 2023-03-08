package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/pkg/shutdown"
	"github.com/indiependente/shrtnr/repository"
	"github.com/indiependente/shrtnr/server"
	"github.com/indiependente/shrtnr/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	appName        = "shrtnr"
	mongoDBTimeout = 5
)

//go:embed ui/dist
var assets embed.FS

func main() {
	err := run()
	if err != nil {
		log.Fatal("unexpected failure: ", err)
	}
}

func run() error { //nolint:funlen
	log := logger.GetLoggerString(appName, "DEBUG")
	conf, err := parseConfig()
	if err != nil {
		return err
	}

	mongoConf := repository.DBConfig{
		User:       conf.MongoDBUser,
		Pass:       conf.MongoDBPassword,
		Host:       conf.MongoDBHost,
		Port:       conf.MongoDBPort,
		DB:         conf.MongoDBName,
		Collection: conf.MongoDBCollection,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConf.URI()))
	if err != nil {
		return fmt.Errorf("could not connect to mongodb: %w", err)
	}
	defer client.Disconnect(ctx) //nolint: errcheck

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("could not ping mongodb: %w", err)
	}
	db := client.Database(mongoConf.DB)

	// create store
	coll := db.Collection(mongoConf.Collection)
	store := repository.NewMongoDBURLStorer(coll)

	// create slugger
	slugger := service.NewFixedLenSlugger(conf.SlugLen)

	// create service
	svc := service.NewURLService(store, slugger)

	// create server
	r := chi.NewRouter()
	srv, err := server.NewHTTPServer(r, svc, conf.Port, http.FS(mustGetFrontend()), log)
	if err != nil {
		return fmt.Errorf("error while creating server: %w", err)
	}

	// // Start HTTP server
	go func() {
		log.Info("Server started listening on " + srv.Addr)
		err := srv.Start(ctx)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("Server closed")

				return
			}
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

func mustGetFrontend() fs.FS {
	f, err := fs.Sub(assets, "ui/dist")
	if err != nil {
		panic(err)
	}

	return f
}
