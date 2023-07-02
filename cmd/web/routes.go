package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxheckel/parks/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/city/clusters", handlers.GetClusters)
}
