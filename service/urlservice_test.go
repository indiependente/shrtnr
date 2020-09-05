package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/repository"
	"github.com/stretchr/testify/require"
)

func TestURLService_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		slugLen           int
		setupExpectations func(storer *repository.MockStorer)
		url               models.URLShortened
		wanterr           bool
	}{
		{
			name:    "Happy Path",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
			},
			url: models.URLShortened{
				URL: "http://indiependente.dev",
			},
			wanterr: false,
		},
		{
			name:              "Sad Path - zero length slug",
			slugLen:           0,
			setupExpectations: func(store *repository.MockStorer) {},
			url: models.URLShortened{
				URL: "http://indiependente.dev",
			},
			wanterr: true,
		},
		{
			name:    "Sad Path - slug in use",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Add(gomock.Any(), gomock.Any()).Return(repository.ErrSlugAlreadyInUse)
			},
			url: models.URLShortened{
				URL: "http://indiependente.dev",
			},
			wanterr: true,
		},
		{
			name:    "Sad Path - unexpected error",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Add(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
			},
			url: models.URLShortened{
				URL: "http://indiependente.dev",
			},
			wanterr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := repository.NewMockStorer(ctrl)
			tt.setupExpectations(mockStore)

			usvc := NewURLService(mockStore, tt.slugLen)

			ctx := context.Background()
			err := usvc.Add(ctx, tt.url)
			require.Equal(t, tt.wanterr, err != nil)
		})
	}
}

func TestURLService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		slug              string
		slugLen           int
		setupExpectations func(storer *repository.MockStorer)
		wanterr           bool
	}{
		{
			name:    "Happy Path",
			slug:    "short",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Delete(gomock.Any(), "short").Return(nil)
			},
			wanterr: false,
		},
		{
			name:              "Sad Path - zero length slug",
			slug:              "",
			slugLen:           0,
			setupExpectations: func(store *repository.MockStorer) {},
			wanterr:           true,
		},
		{
			name:    "Sad Path - slug not found",
			slug:    "short",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Delete(gomock.Any(), "short").Return(repository.ErrSlugNotFound)
			},
			wanterr: true,
		},
		{
			name:    "Sad Path - unexpected error",
			slug:    "short",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Delete(gomock.Any(), "short").Return(errors.New("unexpected error"))
			},
			wanterr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := repository.NewMockStorer(ctrl)
			tt.setupExpectations(mockStore)

			usvc := NewURLService(mockStore, tt.slugLen)

			ctx := context.Background()
			err := usvc.Delete(ctx, tt.slug)
			require.Equal(t, tt.wanterr, err != nil)
		})
	}
}

func TestURLService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		slug              string
		slugLen           int
		setupExpectations func(storer *repository.MockStorer)
		url               models.URLShortened
		wanterr           bool
	}{
		{
			name:    "Happy Path",
			slug:    "short",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Get(gomock.Any(), "short").Return(models.URLShortened{
					Slug: "short",
					URL:  "http://indiependente.dev",
					Hits: 1,
				}, nil)
				store.EXPECT().Update(gomock.Any(), models.URLShortened{
					Slug: "short",
					URL:  "http://indiependente.dev",
					Hits: 2,
				}).MaxTimes(1).Return(nil)
			},
			url: models.URLShortened{
				Slug: "short",
				URL:  "http://indiependente.dev",
				Hits: 1,
			},
			wanterr: false,
		},
		{
			name:              "Sad Path - zero length slug",
			slug:              "",
			slugLen:           0,
			setupExpectations: func(store *repository.MockStorer) {},
			url:               models.URLShortened{},
			wanterr:           true,
		},
		{
			name:    "Sad Path - slug not found",
			slug:    "short",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Get(gomock.Any(), "short").Return(models.URLShortened{}, repository.ErrSlugNotFound)
			},
			url:     models.URLShortened{},
			wanterr: true,
		},
		{
			name:    "Sad Path - unexpected error",
			slug:    "short",
			slugLen: 5,
			setupExpectations: func(store *repository.MockStorer) {
				store.EXPECT().Get(gomock.Any(), "short").Return(models.URLShortened{}, errors.New("unexpected error"))
			},
			url:     models.URLShortened{},
			wanterr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := repository.NewMockStorer(ctrl)
			tt.setupExpectations(mockStore)

			usvc := NewURLService(mockStore, tt.slugLen)

			ctx := context.Background()
			url, err := usvc.Get(ctx, tt.slug)
			require.Equal(t, tt.wanterr, err != nil)
			require.Equal(t, tt.url, url)
		})
	}
}
