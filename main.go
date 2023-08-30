package main

import (
	"govtech-opencv/database"
	"govtech-opencv/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func main() {
		app := fiber.New(fiber.Config{
			// Prefork:       true,
			CaseSensitive: true,
			ServerHeader:  "Fiber",
			AppName: "govtech-opencv",
		})

		app.Use(cors.New())

		database.ConnectToDB()

		router.SetupRoutes(app)

    log.Fatal(app.Listen(":3000"))
}