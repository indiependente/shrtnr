package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/golang/mock/gomock"
	"github.com/indiependente/pkg/logger"
	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/service"
	"github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		slug              string
		setupExpectations func(*service.MockService)
		wantStatus        int
		want              models.URLShortened
	}{
		{
			name: "Happy path",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Get(gomock.Any(), "pizza").Return(models.URLShortened{
					URL:  "http://pizza.com",
					Slug: "pizza",
					Hits: 1000,
				}, nil)
			},
			wantStatus: http.StatusOK,
			want: models.URLShortened{
				URL:  "http://pizza.com",
				Slug: "pizza",
				Hits: 1000,
			},
		},
		{
			name:              "Sad path - empty slug",
			slug:              "",
			setupExpectations: func(mockService *service.MockService) {},
			wantStatus:        http.StatusNotFound,
			want:              models.URLShortened{},
		},
		{
			name: "Sad path - Slug not found",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Get(gomock.Any(), "pizza").Return(models.URLShortened{}, service.ErrSlugNotFound)
			},
			wantStatus: http.StatusNotFound,
			want:       models.URLShortened{},
		},
		{
			name: "Sad path - Slug not valid",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Get(gomock.Any(), "pizza").Return(models.URLShortened{}, service.ErrInvalidSlug)
			},
			wantStatus: http.StatusBadRequest,
			want:       models.URLShortened{},
		},
		{
			name: "Sad path - Unexpected error",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Get(gomock.Any(), "pizza").Return(models.URLShortened{}, errors.New("unexpected error"))
			},
			wantStatus: http.StatusInternalServerError,
			want:       models.URLShortened{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSvc := service.NewMockService(ctrl)
			tt.setupExpectations(mockSvc)

			app := fiber.New(&fiber.Settings{
				CaseSensitive:    true,
				StrictRouting:    true,
				ServerHeader:     "Fiber",
				DisableKeepalive: true, // this is needed to avoid the shutdown being stuck for 30-60 seconds
			})
			srv := NewHTTPServer(app, mockSvc, 9000, logger.GetLogger("test", logger.DISABLED))
			err := srv.Setup(ctx)
			require.NoError(t, err)

			// Start HTTP server
			go func() {
				err := srv.Start(ctx)
				if err != nil {
					t.Error(err)
				}
			}()

			// build request
			path := "http://localhost:9000/url/" + tt.slug
			req, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)
			// send request to server
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close() // nolint: errcheck
			// check status code
			require.Equal(t, tt.wantStatus, resp.StatusCode)
			// parse and check response body
			url := models.URLShortened{}
			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			err = json.Unmarshal(data, &url)
			require.Equal(t, tt.want, url)
			err = srv.Shutdown(ctx) // nolint: errcheck
			require.NoError(t, err)
		})
	}
}

func TestPutURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		url               models.URLShortened
		setupExpectations func(*service.MockService)
		wantStatus        int
		want              models.URLShortened
	}{
		{
			name: "Happy path",
			url: models.URLShortened{
				URL:  "http://pizza.com",
				Slug: "pizza",
			},
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Add(gomock.Any(), models.URLShortened{
					URL:  "http://pizza.com",
					Slug: "pizza",
				}).Return(models.URLShortened{
					URL:  "http://pizza.com",
					Slug: "pizza",
					Hits: 0,
				}, nil)
			},
			wantStatus: http.StatusOK,
			want: models.URLShortened{
				URL:  "http://pizza.com",
				Slug: "pizza",
				Hits: 0,
			},
		},
		{
			name: "Sad path - slug in use",
			url: models.URLShortened{
				URL:  "http://pizza.com",
				Slug: "pizza",
			},
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Add(gomock.Any(), models.URLShortened{
					URL:  "http://pizza.com",
					Slug: "pizza",
				}).Return(models.URLShortened{}, service.ErrSlugAlreadyInUse)
			},
			wantStatus: http.StatusBadRequest,
			want:       models.URLShortened{},
		},
		{
			name: "Sad path - Slug not valid",
			url: models.URLShortened{
				URL:  "http://pizza.com",
				Slug: "pizza",
			},
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Add(gomock.Any(), models.URLShortened{
					URL:  "http://pizza.com",
					Slug: "pizza",
				}).Return(models.URLShortened{}, service.ErrInvalidSlug)
			},
			wantStatus: http.StatusBadRequest,
			want:       models.URLShortened{},
		},
		{
			name: "Sad path - Unexpected error",
			url: models.URLShortened{
				URL:  "http://pizza.com",
				Slug: "pizza",
			},
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Add(gomock.Any(), models.URLShortened{
					URL:  "http://pizza.com",
					Slug: "pizza",
				}).Return(models.URLShortened{}, errors.New("unexpected error"))
			},
			wantStatus: http.StatusInternalServerError,
			want:       models.URLShortened{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSvc := service.NewMockService(ctrl)
			tt.setupExpectations(mockSvc)

			app := fiber.New(&fiber.Settings{
				CaseSensitive:    true,
				StrictRouting:    true,
				ServerHeader:     "Fiber",
				DisableKeepalive: true, // this is needed to avoid the shutdown being stuck for 30-60 seconds
			})
			srv := NewHTTPServer(app, mockSvc, 9000, logger.GetLogger("test", logger.DISABLED))
			err := srv.Setup(ctx)
			require.NoError(t, err)

			// Start HTTP server
			go func() {
				err := srv.Start(ctx)
				if err != nil {
					t.Error(err)
				}
			}()

			// build request
			path := "http://localhost:9000/url"
			reqBody, err := json.Marshal(tt.url)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			// send request to server
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close() // nolint: errcheck
			// check status code
			require.Equal(t, tt.wantStatus, resp.StatusCode)
			// parse and check response body
			url := models.URLShortened{}
			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			err = json.Unmarshal(data, &url)
			require.Equal(t, tt.want, url)
			err = srv.Shutdown(ctx) // nolint: errcheck
			require.NoError(t, err)
		})
	}
}

func TestDeleteURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		slug              string
		setupExpectations func(*service.MockService)
		wantStatus        int
	}{
		{
			name: "Happy path",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Delete(gomock.Any(), "pizza").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:              "Sad path - empty slug",
			slug:              "",
			setupExpectations: func(mockService *service.MockService) {},
			wantStatus:        http.StatusNotFound,
		},
		{
			name: "Sad path - Slug not found",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Delete(gomock.Any(), "pizza").Return(service.ErrSlugNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Sad path - Slug not valid",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Delete(gomock.Any(), "pizza").Return(service.ErrInvalidSlug)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Sad path - Unexpected error",
			slug: "pizza",
			setupExpectations: func(mockService *service.MockService) {
				mockService.EXPECT().Delete(gomock.Any(), "pizza").Return(errors.New("unexpected error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSvc := service.NewMockService(ctrl)
			tt.setupExpectations(mockSvc)

			app := fiber.New(&fiber.Settings{
				CaseSensitive:    true,
				StrictRouting:    true,
				ServerHeader:     "Fiber",
				DisableKeepalive: true, // this is needed to avoid the shutdown being stuck for 30-60 seconds
			})
			srv := NewHTTPServer(app, mockSvc, 9000, logger.GetLogger("test", logger.DISABLED))
			err := srv.Setup(ctx)
			require.NoError(t, err)

			// Start HTTP server
			go func() {
				err := srv.Start(ctx)
				if err != nil {
					t.Error(err)
				}
			}()

			// build request
			path := "http://localhost:9000/url/" + tt.slug
			req, err := http.NewRequest(http.MethodDelete, path, nil)
			require.NoError(t, err)
			// send request to server
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close() // nolint: errcheck
			// check status code
			require.Equal(t, tt.wantStatus, resp.StatusCode)

			err = srv.Shutdown(ctx) // nolint: errcheck
			require.NoError(t, err)
		})
	}
}