package server

import (
	"net/http"

	"github.com/gofiber/fiber"
	"github.com/indiependente/shrtnr/models"
	"github.com/indiependente/shrtnr/service"
)

func getURL(service service.Service) fiber.Handler {
	return func(c *fiber.Ctx) {
		slug := c.Params("slug")
		url, err := service.Get(c.Context(), slug)
		if err != nil {
			c.SendStatus(http.StatusBadRequest)
			return
		}
		if err := c.Status(http.StatusOK).JSON(url); err != nil {
			c.Status(http.StatusInternalServerError).Send(err)
			return
		}
	}
}

func putURL(service service.Service) fiber.Handler {
	return func(c *fiber.Ctx) {
		url := models.URLShortened{}
		if err := c.BodyParser(&url); err != nil {
			c.SendStatus(http.StatusBadRequest)
			return
		}
		newUrl, err := service.Add(c.Context(), url)
		if err != nil {
			c.Status(500).Send(err)
			return
		}
		if err := c.Status(http.StatusAccepted).JSON(newUrl); err != nil {
			c.Status(http.StatusInternalServerError).Send(err)
			return
		}
	}
}

func delURL(service service.Service) fiber.Handler {
	return func(c *fiber.Ctx) {
		slug := c.Params("slug")
		err := service.Delete(c.Context(), slug)
		if err != nil {
			c.Status(http.StatusInternalServerError).Send(err)
			return
		}
		c.SendStatus(http.StatusOK)
	}
}
