package service

import (
	"context"
	"fmt"

	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/repository"
)

// URLService implements the Service interface.
type URLService struct {
	store   repository.Storer
	slugLen int
}

// NewURLService returns a new instance of the URLService type.
func NewURLService(store repository.Storer, slugLen int) URLService {
	return URLService{
		store:   store,
		slugLen: slugLen,
	}
}

func (usvc URLService) Add(ctx context.Context, shortURL models.URLShortened) error {
	if shortURL.Slug == "" {
		shortURL.Slug = generateSlug(usvc.slugLen)
	}
	if !validateSlug(shortURL.Slug, usvc.slugLen) {
		return fmt.Errorf("could not use slug: %w", ErrInvalidSlug)
	}
	err := usvc.store.Add(ctx, shortURL)
	if err != nil {
		if err == repository.ErrSlugAlreadyInUse {
			return fmt.Errorf("could not add: %w", ErrSlugAlreadyInUse)
		}
		return fmt.Errorf("could not add: %w", err)
	}
	return nil
}

func (usvc URLService) Get(ctx context.Context, slug string) (models.URLShortened, error) {
	if slug == "" {
		return models.URLShortened{}, fmt.Errorf("empty slug: %w", ErrInvalidSlug)
	}
	url, err := usvc.store.Get(ctx, slug)
	if err != nil {
		if err == repository.ErrSlugNotFound {
			return models.URLShortened{}, fmt.Errorf("could not get: %w", ErrSlugNotFound)
		}
		return models.URLShortened{}, fmt.Errorf("could not get: %w", err)
	}
	// increase hit counter and store the updated value
	go usvc.increaseHitCounter(url)
	return url, nil
}

// increaseHitCounter increases the hit count by one and updates the value in the repo.
// It is supposed to be called in a separate goroutine.
func (usvc URLService) increaseHitCounter(url models.URLShortened) {
	url.Hits++
	_ = usvc.store.Update(context.Background(), url)
}

// Delete deletes the entry related to the input slug from the repository.
func (usvc URLService) Delete(ctx context.Context, slug string) error {
	if slug == "" {
		return fmt.Errorf("empty slug: %w", ErrInvalidSlug)
	}
	err := usvc.store.Delete(ctx, slug)
	if err != nil {
		if err == repository.ErrSlugNotFound {
			return fmt.Errorf("could not delete: %w", ErrSlugNotFound)
		}
		return fmt.Errorf("could not delete: %w", err)
	}
	return nil
}
