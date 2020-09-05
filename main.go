package main

import (
	"log"

	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, World 👋!")
	})

	log.Fatal(app.Listen(7000))
}
