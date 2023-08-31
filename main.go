package main

import (
	"govtech-opencv/db"
	"govtech-opencv/router"
	"log"

	_ "govtech-opencv/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		ServerHeader:  "Fiber",
		AppName:       "govtech-opencv",
	})

	app.Use(cors.New())
	// app.Get("/swagger/*", swagger.HandlerDefault)

	db.ConnectToDB()

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
