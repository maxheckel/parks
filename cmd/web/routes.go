package main

import (
	"encoding/gob"
	"github.com/gofiber/fiber/v2"
	"github.com/maxheckel/parks/handlers"
	"github.com/maxheckel/parks/handlers/auth"
)

func setupRoutes(app *fiber.App) {
	gob.Register(map[string]interface{}{})
	app.Get("/city/clusters", handlers.GetClusters)
	app.Get("/auth/callback", auth.Callback)
	app.Get("/auth/login", auth.Login)
	app.Get("/auth/logout", auth.Logout)
	app.Static("/", "./web/dist")

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./web/dist/index.html")
	})
}
