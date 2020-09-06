package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/indiependente/shrtnr/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	uriFmt = "mongodb://%s:%s@%s:%s/%s"
)

// mongoURLShortened is the model representation of the data for the mongo database.
type mongoURLShortened struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	URL  string             `bson:"url"`
	Slug string             `bson:"slug"`
	Hits int                `bson:"hits"`
}

// MongoDBStorer implements the Storer using a MongoDB store.
type MongoDBURLStorer struct {
	urls *mongo.Collection
}

// NewMongoDBURLStorer returns a new instance of a MongoDBURLStorer.
func NewMongoDBURLStorer(coll *mongo.Collection) MongoDBURLStorer {
	return MongoDBURLStorer{
		urls: coll,
	}
}

// Add adds a shortened url to the mongodb repository.
// Returns an error if any.
func (m MongoDBURLStorer) Add(ctx context.Context, shortened models.URLShortened) error {
	url, err := m.Get(ctx, shortened.Slug)
	if err != nil && !errors.Is(err, ErrSlugNotFound) { // err can be ErrSlugNotFound
		return fmt.Errorf("could not lookup: %w", err)
	}
	if url.Slug == shortened.Slug {
		return fmt.Errorf("could not add: %w", ErrSlugAlreadyInUse)
	}
	_, err = m.urls.InsertOne(ctx, toMongo(shortened))
	if err != nil {
		return fmt.Errorf("could not insert: %w", err)
	}
	return nil
}

// Get gets a shortened url from the mongodb repository.
// Returns an error if any.
func (m MongoDBURLStorer) Get(ctx context.Context, slug string) (models.URLShortened, error) {
	var shortURL mongoURLShortened
	err := m.urls.FindOne(ctx, bson.D{{Key: "slug", Value: slug}}).Decode(&shortURL)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.URLShortened{}, ErrSlugNotFound
		}
		return models.URLShortened{}, fmt.Errorf("unexpected error: %w", err)
	}
	return toModel(shortURL), nil
}

// Update deletes a shortened url from the mongodb repository.
// Returns an error if any.
func (m MongoDBURLStorer) Update(ctx context.Context, newshort models.URLShortened) error {
	filter := bson.D{{Key: "slug", Value: newshort.Slug}}
	update := bson.D{
		{Key: "$set", Value: toMongo(newshort)},
		{Key: "$currentDate", Value: bson.D{
			{Key: "lastModified", Value: true},
		}},
	}
	result, err := m.urls.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("could not update: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("could not update: %w", ErrSlugNotFound)
	}
	return nil
}

// Delete deletes a shortened url from the mongodb repository.
// Returns an error if any.
func (m MongoDBURLStorer) Delete(ctx context.Context, slug string) error {
	res, err := m.urls.DeleteOne(ctx, bson.D{{Key: "slug", Value: slug}})
	if err != nil {
		return fmt.Errorf("could not delete: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("could not delete: %w", ErrSlugNotFound)
	}
	return nil
}

// Configs MongoDB configuration
type Configs struct {
	User, Pass, Host, Port, DB, Collection string
}

// URI returns the URI string.
func (c Configs) URI() string {
	return fmt.Sprintf(uriFmt, c.User, c.Pass, c.Host, c.Port, c.DB)
}

// BuildMongoConfigs parses MongoDB configurations from the environment.
func BuildMongoConfigs() Configs {
	user := os.Getenv("MONGODB_USER")
	pass := os.Getenv("MONGODB_PASSWORD")
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")
	db := os.Getenv("MONGODB_DB")
	coll := os.Getenv("MONGODB_COLLECTION")
	return Configs{User: user, Pass: pass, Host: host, Port: port, DB: db, Collection: coll}
}

func toMongo(u models.URLShortened) mongoURLShortened {
	return mongoURLShortened{
		URL:  u.URL,
		Slug: u.Slug,
		Hits: u.Hits,
	}
}

func toModel(mu mongoURLShortened) models.URLShortened {
	return models.URLShortened{
		URL:  mu.URL,
		Slug: mu.Slug,
		Hits: mu.Hits,
	}
}
