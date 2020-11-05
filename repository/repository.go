//go:generate mockgen -package repository -source=repository.go -destination repository_mock.go

package repository

import (
	"context"

	"github.com/indiependente/shrtnr/models"
)

const (
	// ErrSlugAlreadyInUse is returned when trying to store a shortened url with a slug already in use.
	ErrSlugAlreadyInUse Error = `slug in use`
	// ErrSlugNotFound is returned when trying to retrieve a slug that could not be found in the repository.
	ErrSlugNotFound Error = `slug not found`
	// ErrURLNotFound is returned when trying to retrieve a URL that could not be found in the repository.
	ErrURLNotFound Error = `url not found`
)

// Error represents an error returned by the repository.
type Error string

// Error returns the string representation of the error.
func (e Error) Error() string {
	return string(e)
}

// Storer defines the behaviour of a component capable of storing shortened urls, retrieving and deleting existing ones.
type Storer interface {
	Add(ctx context.Context, shortened models.URLShortened) error
	Get(ctx context.Context, slug string) (models.URLShortened, error)
	GetURL(ctx context.Context, url string) (models.URLShortened, error)
	Update(ctx context.Context, newshortened models.URLShortened) error
	Delete(ctx context.Context, slug string) error
}
