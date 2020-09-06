package server

import (
	"net/http"

	"github.com/gofiber/fiber"
	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/service"
)

func getURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) {
		slug := c.Params("slug")
		url, err := svc.Get(c.Context(), slug)
		switch {
		case err == service.ErrSlugNotFound:
			c.SendStatus(http.StatusNotFound)
			return
		case err == service.ErrInvalidSlug:
			c.SendStatus(http.StatusBadRequest)
			return
		case err != nil:
			c.Status(http.StatusInternalServerError).Send(err)
			return
		default: // all good
			if err := c.Status(http.StatusOK).JSON(url); err != nil {
				c.Status(http.StatusInternalServerError).Send(err)
				return
			}
		}
	}
}

func putURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) {
		url := models.URLShortened{}
		if err := c.BodyParser(&url); err != nil {
			c.Status(http.StatusBadRequest).Send(err)
			return
		}
		newUrl, err := svc.Add(c.Context(), url)
		switch {
		case err == service.ErrSlugAlreadyInUse:
			c.Status(http.StatusBadRequest).Send(err)
			return
		case err == service.ErrInvalidSlug:
			c.Status(http.StatusBadRequest).Send(err)
			return
		case err != nil:
			c.Status(http.StatusInternalServerError).Send(err)
			return
		default: // all good
			if err := c.Status(http.StatusOK).JSON(newUrl); err != nil {
				c.Status(http.StatusInternalServerError).Send(err)
				return
			}
		}
	}
}

func delURL(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) {
		slug := c.Params("slug")
		err := svc.Delete(c.Context(), slug)
		switch {
		case err == service.ErrSlugNotFound:
			c.SendStatus(http.StatusNotFound)
			return
		case err == service.ErrInvalidSlug:
			c.SendStatus(http.StatusBadRequest)
			return
		case err != nil:
			c.Status(http.StatusInternalServerError).Send(err)
			return
		default: // all good
			c.SendStatus(http.StatusOK)
		}
	}
}
