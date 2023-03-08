package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/service"
)

func (srv *HTTPServer) getURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		url, err := srv.svc.Get(r.Context(), slug)
		switch {
		case errors.Is(err, service.ErrSlugNotFound):
			http.Error(w, "not found", http.StatusNotFound)

			return
		case errors.Is(err, service.ErrInvalidSlug):
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		case err != nil:
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		default: // all good
			err := json.NewEncoder(w).Encode(url)
			if err != nil {
				http.Error(w, "encoding error", http.StatusInternalServerError)

				return
			}
		}
	}
}

func (srv *HTTPServer) putURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := models.URLShortened{}
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		newURL, err := srv.svc.Add(r.Context(), url)
		switch {
		case errors.Is(err, service.ErrSlugAlreadyInUse):
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		case errors.Is(err, service.ErrInvalidSlug):
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		case err != nil:
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		default: // all good
			err := json.NewEncoder(w).Encode(newURL)
			if err != nil {
				http.Error(w, "encoding error", http.StatusInternalServerError)

				return
			}
		}
	}
}

func (srv *HTTPServer) delURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		err := srv.svc.Delete(r.Context(), slug)
		switch {
		case errors.Is(err, service.ErrSlugNotFound):
			http.Error(w, "bad request", http.StatusNotFound)

			return
		case errors.Is(err, service.ErrInvalidSlug):
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		case err != nil:
			http.Error(w, "encoding error", http.StatusInternalServerError)

			return
		default: // all good
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (srv *HTTPServer) resolveURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		url, err := srv.svc.Get(r.Context(), slug)
		switch {
		case errors.Is(err, service.ErrSlugNotFound):
			http.Error(w, "not found", http.StatusNotFound)

			return
		case errors.Is(err, service.ErrInvalidSlug):
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		case err != nil:
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		default: // all good
			http.Redirect(w, r, url.URL, http.StatusMovedPermanently)

			return
		}
	}
}

func (srv *HTTPServer) shortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := models.URLShortened{}
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		short, err := srv.svc.Shorten(r.Context(), url.URL)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}
		err = json.NewEncoder(w).Encode(short)
		if err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)

			return
		}
	}
}
