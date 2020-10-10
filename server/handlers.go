package server

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/service"
)

func getURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		url, err := svc.Get(c.Context(), slug)
		switch {
		case errors.Is(err, service.ErrSlugNotFound):
			return c.SendStatus(http.StatusNotFound)
		case errors.Is(err, service.ErrInvalidSlug):
			return c.SendStatus(http.StatusBadRequest)
		case err != nil:
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		default: // all good
			if err := c.Status(http.StatusOK).JSON(url); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
		}
		return nil
	}
}

func putURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		url := models.URLShortened{}
		if err := c.BodyParser(&url); err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		newUrl, err := svc.Add(c.Context(), url)
		switch {
		case errors.Is(err, service.ErrSlugAlreadyInUse):
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		case errors.Is(err, service.ErrInvalidSlug):
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		case err != nil:
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		default: // all good
			if err := c.Status(http.StatusOK).JSON(newUrl); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
		}
		return nil
	}
}

func delURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		err := svc.Delete(c.Context(), slug)
		switch {
		case errors.Is(err, service.ErrSlugNotFound):
			return c.SendStatus(http.StatusNotFound)
		case errors.Is(err, service.ErrInvalidSlug):
			return c.SendStatus(http.StatusBadRequest)
		case err != nil:
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		default: // all good
			return c.SendStatus(http.StatusOK)
		}
	}
}

func resolveURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		url, err := svc.Get(c.Context(), slug)
		switch {
		case errors.Is(err, service.ErrSlugNotFound):
			return c.SendStatus(http.StatusNotFound)
		case errors.Is(err, service.ErrInvalidSlug):
			return c.SendStatus(http.StatusBadRequest)
		case err != nil:
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		default: // all good
			return c.Redirect(url.URL, http.StatusMovedPermanently)
		}
	}
}
