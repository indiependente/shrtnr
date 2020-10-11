package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/repository"
)

// URLService implements the Service interface.
type URLService struct {
	store   repository.Storer
	slugger Slugger
}

// NewURLService returns a new instance of the URLService type.
func NewURLService(store repository.Storer, slugger Slugger) URLService {
	return URLService{
		store:   store,
		slugger: slugger,
	}
}

func (usvc URLService) Add(ctx context.Context, shortURL models.URLShortened) (models.URLShortened, error) {
	if shortURL.Slug == "" {
		shortURL.Slug = usvc.slugger.Slug()
	}
	if !usvc.slugger.Validate(shortURL.Slug) {
		return models.URLShortened{}, fmt.Errorf("could not use slug: %w", ErrInvalidSlug)
	}
	err := usvc.store.Add(ctx, shortURL)
	if err != nil {
		if errors.Is(err, repository.ErrSlugAlreadyInUse) {
			return models.URLShortened{}, fmt.Errorf("could not add: %w", ErrSlugAlreadyInUse)
		}
		return models.URLShortened{}, fmt.Errorf("could not add: %w", err)
	}
	return shortURL, nil
}

func (usvc URLService) Get(ctx context.Context, slug string) (models.URLShortened, error) {
	if slug == "" {
		return models.URLShortened{}, fmt.Errorf("empty slug: %w", ErrInvalidSlug)
	}
	url, err := usvc.store.Get(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrSlugNotFound) {
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

// Shorten returns the shortened URL and shortens it if not found.
// Returns an error if any.
func (usvc URLService) Shorten(ctx context.Context, url string) (models.URLShortened, error) {
	url = fixURL(url)
	short, err := usvc.store.GetURL(ctx, url) // try to get from repo
	if err != nil {
		if !errors.Is(err, repository.ErrURLNotFound) {
			return models.URLShortened{}, fmt.Errorf("could not shorten: %w", err)
		}
		// create if not found
		short.URL = url
		short.Slug = usvc.slugger.Slug()
		err := usvc.store.Add(ctx, short)
		if err != nil {
			if errors.Is(err, repository.ErrSlugAlreadyInUse) {
				return models.URLShortened{}, fmt.Errorf("could not add: %w", ErrSlugAlreadyInUse)
			}
			return models.URLShortened{}, fmt.Errorf("could not add: %w", err)
		}
	}
	// increase hit counter and store the updated value
	go usvc.increaseHitCounter(short)
	return short, nil
}

// Delete deletes the entry related to the input slug from the repository.
func (usvc URLService) Delete(ctx context.Context, slug string) error {
	if slug == "" {
		return fmt.Errorf("empty slug: %w", ErrInvalidSlug)
	}
	err := usvc.store.Delete(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrSlugNotFound) {
			return fmt.Errorf("could not delete: %w", ErrSlugNotFound)
		}
		return fmt.Errorf("could not delete: %w", err)
	}
	return nil
}

func fixURL(url string) string {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	return url
}
