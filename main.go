package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/pkg/shutdown"
	"github.com/indiependente/shrtnr/repository"
	"github.com/indiependente/shrtnr/server"
	"github.com/indiependente/shrtnr/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	appName        = "shrtnr"
	mongoDBTimeout = 2
)

//go:embed ui/dist
var assets embed.FS

func main() {
	err := run()
	if err != nil {
		log.Fatal("unexpected failure: ", err)
	}
}

func run() error { //nolint:funlen,cyclop
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
	ctx, cancel := context.WithTimeout(context.Background(), mongoDBTimeout*time.Second)
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
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	var results []any
	err = cur.All(ctx, &results)
	if err != nil {
		return err
	}
	for _, u := range results {
		fmt.Println(u)
	}
	store := repository.NewMongoDBURLStorer(coll)

	// create slugger
	slugger := service.NewFixedLenSlugger(conf.SlugLen)

	// create service
	svc := service.NewURLService(store, slugger)

	// create server
	r := chi.NewRouter()
	srv, err := server.NewHTTPServer(r, svc, conf.Port, http.FS(assets), log)
	if err != nil {
		return fmt.Errorf("error while creating server: %w", err)
	}
	err = srv.Setup(ctx)
	if err != nil {
		return fmt.Errorf("error while running server setup: %w", err)
	}

	// Start HTTP server
	go func() {
		err := srv.Start(ctx) //nolint:govet
		if err != nil {
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
