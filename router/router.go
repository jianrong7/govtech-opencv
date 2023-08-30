package router

import (
	"govtech-opencv/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})
	
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)
}