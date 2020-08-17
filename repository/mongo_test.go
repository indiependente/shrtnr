package repository

import (
	"context"
	"errors"
	"fmt"
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
			// mongodb setup
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			err = client.Connect(ctx)
			require.NoError(t, err)
			defer client.Disconnect(ctx)
			db := client.Database("shrtnr")
			// create collection
			coll := db.Collection(fmt.Sprintf("urls_test_add_%d", time.Now().UnixNano()))
			defer coll.Drop(ctx)
			// add indexes
			models := []mongo.IndexModel{
				{
					Keys: bson.D{{"slug", 1}},
				},
				{
					Keys: bson.D{{"url", 1}},
				},
			}
			opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
			_, err = coll.Indexes().CreateMany(ctx, models, opts)
			require.NoError(t, err)
			defer coll.Indexes().DropAll(ctx)
			// run additional collection setup func
			err = tt.setupCollection(ctx, coll)
			require.NoError(t, err)
			// create store and test Add
			store := NewMongoDBURLStorer(coll)
			err = store.Add(ctx, tt.url)
			require.True(t, errors.Is(err, tt.err))

		})
	}
}

//
//func TestMongoDBURLStorer_Delete(t *testing.T) {
//	type fields struct {
//		urls mongo.Collection
//	}
//	type args struct {
//		ctx  context.Context
//		slug string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := MongoDBURLStorer{
//				urls: tt.fields.urls,
//			}
//			if err := m.Delete(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
//				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestMongoDBURLStorer_Get(t *testing.T) {
//	type fields struct {
//		urls mongo.Collection
//	}
//	type args struct {
//		ctx  context.Context
//		slug string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    models.URLShortened
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := MongoDBURLStorer{
//				urls: tt.fields.urls,
//			}
//			got, err := m.Get(tt.args.ctx, tt.args.slug)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Get() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMongoDBURLStorer_Update(t *testing.T) {
//	type fields struct {
//		urls mongo.Collection
//	}
//	type args struct {
//		ctx      context.Context
//		newshort models.URLShortened
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := MongoDBURLStorer{
//				urls: tt.fields.urls,
//			}
//			if err := m.Update(tt.args.ctx, tt.args.newshort); (err != nil) != tt.wantErr {
//				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestNewMongoDBURLStorer(t *testing.T) {
//	type args struct {
//		coll mongo.Collection
//	}
//	tests := []struct {
//		name string
//		args args
//		want MongoDBURLStorer
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewMongoDBURLStorer(tt.args.coll); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewMongoDBURLStorer() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_toModel(t *testing.T) {
//	type args struct {
//		mu mongoURLShortened
//	}
//	tests := []struct {
//		name string
//		args args
//		want models.URLShortened
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := toModel(tt.args.mu); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("toModel() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_toMongo(t *testing.T) {
//	type args struct {
//		u models.URLShortened
//	}
//	tests := []struct {
//		name string
//		args args
//		want mongoURLShortened
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := toMongo(tt.args.u); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("toMongo() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
