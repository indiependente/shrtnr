//go:generate mockgen -package service -source=service.go -destination service_mock.go

package service

import (
	"context"

	"github.com/indiependente/shrtnr/models"
)

const (
	// ErrSlugAlreadyInUse is returned when trying to add a shortened url with a slug already in use.
	ErrSlugAlreadyInUse Error = `slug in use`
	// ErrSlugNotFound is returned when trying to get or delete a slug that could not be found in the service.
	ErrSlugNotFound Error = `slug not found`
	// ErrInvalidSlug is returned when trying to use a not valid slug.
	ErrInvalidSlug Error = `slug not valid`
)

// Error represents an error returned by the repository.
type Error string

// Error returns the string representation of the error.
func (e Error) Error() string {
	return string(e)
}

// Service defines the behaviour of a service capable of shortening urls, retrieving and deleting shortened ones.
type Service interface {
	Add(ctx context.Context, shortURL models.URLShortened) (models.URLShortened, error)
	Get(ctx context.Context, slug string) (models.URLShortened, error)
	Delete(ctx context.Context, slug string) error
}
