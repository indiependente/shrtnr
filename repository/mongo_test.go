//go:build integration
// +build integration

package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/indiependente/shrtnr/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri = "mongodb://frank:password@localhost:27017/shrtnr"
)

func TestMongoDBURLStorer_Add(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		url             models.URLShortened
		setupCollection func(ctx context.Context, coll *mongo.Collection) error
		err             error
	}{
		{
			name: "Happy path",
			url: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				return nil
			},
			err: nil,
		},
		{
			name: "Sad path - existing slug",
			url: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				_, err := coll.InsertOne(ctx, toMongo(models.URLShortened{
					URL:  "https://shrtnr.dev",
					Slug: "aeiou",
					Hits: 0,
				}))
				return err
			},
			err: ErrSlugAlreadyInUse,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// *** START DB SETUP ***
			rand.Seed(time.Now().UnixNano())
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = client.Connect(ctx)
			require.NoError(t, err)
			defer client.Disconnect(ctx) // nolint: errcheck
			db := client.Database("shrtnr")
			// create collection
			coll := db.Collection(fmt.Sprintf("urls_test_add_%d%d", time.Now().UnixNano(), rand.Int()))
			defer coll.Drop(ctx) // nolint: errcheck
			// add indexes
			models := []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "slug", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "url", Value: 1}},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
			_, err = coll.Indexes().CreateMany(ctx, models, opts)
			require.NoError(t, err)
			defer coll.Indexes().DropAll(ctx) // nolint: errcheck
			// run additional collection setup func
			err = tt.setupCollection(ctx, coll)
			require.NoError(t, err)
			// *** END DB SETUP ***
			// create store and test Add
			store := NewMongoDBURLStorer(coll)
			// add a new url
			err = store.Add(ctx, tt.url)
			require.True(t, errors.Is(err, tt.err))
		})
	}
}

func TestMongoDBURLStorer_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		url             models.URLShortened
		setupCollection func(ctx context.Context, coll *mongo.Collection) error
		err             error
	}{
		{
			name: "Happy path",
			url: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				_, err := coll.InsertOne(ctx, toMongo(models.URLShortened{
					URL:  "https://shrtnr.dev",
					Slug: "aeiou",
					Hits: 0,
				}))
				return err
			},
			err: nil,
		},
		{
			name: "Sad path - slug not found",
			url: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				return nil
			},
			err: ErrSlugNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// *** START DB SETUP ***
			rand.Seed(time.Now().UnixNano())
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = client.Connect(ctx)
			require.NoError(t, err)
			defer client.Disconnect(ctx) // nolint: errcheck
			db := client.Database("shrtnr")
			// create collection
			coll := db.Collection(fmt.Sprintf("urls_test_delete_%d%d", time.Now().UnixNano(), rand.Int()))
			defer coll.Drop(ctx) // nolint: errcheck
			// add indexes
			models := []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "slug", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "url", Value: 1}},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
			_, err = coll.Indexes().CreateMany(ctx, models, opts)
			require.NoError(t, err)
			defer coll.Indexes().DropAll(ctx) // nolint: errcheck
			// run additional collection setup func
			err = tt.setupCollection(ctx, coll)
			require.NoError(t, err)
			// *** END DB SETUP ***
			// create store and test Delete
			store := NewMongoDBURLStorer(coll)
			// delete url
			err = store.Delete(ctx, tt.url.Slug)
			require.True(t, errors.Is(err, tt.err))
		})
	}
}
func TestMongoDBURLStorer_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		slug            string
		setupCollection func(ctx context.Context, coll *mongo.Collection) error
		wantURL         models.URLShortened
		err             error
	}{
		{
			name: "Happy path",
			slug: "aeiou",
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				_, err := coll.InsertOne(ctx, toMongo(models.URLShortened{
					URL:  "https://shrtnr.dev",
					Slug: "aeiou",
					Hits: 0,
				}))
				return err
			},
			wantURL: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			err: nil,
		},
		{
			name: "Sad path - slug not found",
			slug: "aeiou",
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				return nil
			},
			wantURL: models.URLShortened{},
			err:     ErrSlugNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// *** START DB SETUP ***
			rand.Seed(time.Now().UnixNano())
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = client.Connect(ctx)
			require.NoError(t, err)
			defer client.Disconnect(ctx) // nolint: errcheck
			db := client.Database("shrtnr")
			// create collection
			coll := db.Collection(fmt.Sprintf("urls_test_get_%d%d", time.Now().UnixNano(), rand.Int()))
			defer coll.Drop(ctx) // nolint: errcheck
			// add indexes
			models := []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "slug", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "url", Value: 1}},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
			_, err = coll.Indexes().CreateMany(ctx, models, opts)
			require.NoError(t, err)
			defer coll.Indexes().DropAll(ctx) // nolint: errcheck
			// run additional collection setup func
			err = tt.setupCollection(ctx, coll)
			require.NoError(t, err)
			// *** END DB SETUP ***
			// create store and test Get
			store := NewMongoDBURLStorer(coll)
			// delete url
			url, err := store.Get(ctx, tt.slug)
			require.True(t, errors.Is(err, tt.err))
			require.Equal(t, tt.wantURL, url)
		})
	}
}

func TestMongoDBURLStorer_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		url             models.URLShortened
		setupCollection func(ctx context.Context, coll *mongo.Collection) error
		err             error
	}{
		{
			name: "Happy path",
			url: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 1,
			},
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				_, err := coll.InsertOne(ctx, toMongo(models.URLShortened{
					URL:  "https://shrtnr.dev",
					Slug: "aeiou",
					Hits: 0,
				}))
				return err
			},
			err: nil,
		},
		{
			name: "Sad path - slug not found",
			url: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				return nil
			},
			err: ErrSlugNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// *** START DB SETUP ***
			rand.Seed(time.Now().UnixNano())
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = client.Connect(ctx)
			require.NoError(t, err)
			defer client.Disconnect(ctx) // nolint: errcheck
			db := client.Database("shrtnr")
			// create collection
			coll := db.Collection(fmt.Sprintf("urls_test_update_%d%d", time.Now().UnixNano(), rand.Int()))
			defer coll.Drop(ctx) // nolint: errcheck
			// add indexes
			models := []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "slug", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "url", Value: 1}},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
			_, err = coll.Indexes().CreateMany(ctx, models, opts)
			require.NoError(t, err)
			defer coll.Indexes().DropAll(ctx) // nolint: errcheck
			// run additional collection setup func
			err = tt.setupCollection(ctx, coll)
			require.NoError(t, err)
			// *** END DB SETUP ***
			// create store and test Update
			store := NewMongoDBURLStorer(coll)
			// delete url
			err = store.Update(ctx, tt.url)
			require.True(t, errors.Is(err, tt.err))
		})
	}
}

func TestMongoDBURLStorer_GetURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		url             string
		setupCollection func(ctx context.Context, coll *mongo.Collection) error
		wantURL         models.URLShortened
		err             error
	}{
		{
			name: "Happy path",
			url:  "https://shrtnr.dev",
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				_, err := coll.InsertOne(ctx, toMongo(models.URLShortened{
					URL:  "https://shrtnr.dev",
					Slug: "aeiou",
					Hits: 0,
				}))
				return err
			},
			wantURL: models.URLShortened{
				URL:  "https://shrtnr.dev",
				Slug: "aeiou",
				Hits: 0,
			},
			err: nil,
		},
		{
			name: "Sad path - slug not found",
			url:  "https://shrtnr.dev",
			setupCollection: func(ctx context.Context, coll *mongo.Collection) error {
				return nil
			},
			wantURL: models.URLShortened{},
			err:     ErrURLNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// *** START DB SETUP ***
			rand.Seed(time.Now().UnixNano())
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = client.Connect(ctx)
			require.NoError(t, err)
			defer client.Disconnect(ctx) // nolint: errcheck
			db := client.Database("shrtnr")
			// create collection
			coll := db.Collection(fmt.Sprintf("urls_test_geturl_%d%d", time.Now().UnixNano(), rand.Int()))
			defer coll.Drop(ctx) // nolint: errcheck
			// add indexes
			models := []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "slug", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "url", Value: 1}},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
			_, err = coll.Indexes().CreateMany(ctx, models, opts)
			require.NoError(t, err)
			defer coll.Indexes().DropAll(ctx) // nolint: errcheck
			// run additional collection setup func
			err = tt.setupCollection(ctx, coll)
			require.NoError(t, err)
			// *** END DB SETUP ***
			// create store and test Get
			store := NewMongoDBURLStorer(coll)
			// delete url
			url, err := store.GetURL(ctx, tt.url)
			require.True(t, errors.Is(err, tt.err))
			require.Equal(t, tt.wantURL, url)
		})
	}
}
