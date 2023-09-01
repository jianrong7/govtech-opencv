package router

import (
	"govtech-opencv/app/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ping")
	})

	api := app.Group("/api", logger.New())
	api.Post("/register", controller.Register)
	api.Get("/commonstudents", controller.GetCommonStudents)
	api.Post("/suspend", controller.SuspendStudent)
	api.Post("/retrievefornotifications", controller.RetrieveForNotifications)
}
