package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jiteshchawla1511/url-shortener/routes"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.Resolve)
	app.Post("/api/v1", routes.Shorten)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
