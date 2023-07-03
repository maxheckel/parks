package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/maxheckel/parks/handlers"
	"github.com/maxheckel/parks/handlers/auth"
	"github.com/maxheckel/parks/services/store"
)

func setupRoutes(app *fiber.App) {
	gob.Register(map[string]interface{}{})
	app.Get("/city/clusters", handlers.GetClusters)
	app.Get("/auth/callback", auth.Callback)
	app.Get("/auth/login", auth.Login)
	app.Get("/auth/logout", auth.Logout)

	authed := app.Group("/account", Authenticated)
	authed.Get("/user", func(ctx *fiber.Ctx) error {
		session, _ := store.Store.Get(ctx)
		profile := session.Get("profile")
		fmt.Println(profile)
		return ctx.Status(200).JSON(profile)
	})
	app.Static("/", "./web/dist")

	app.Get("/about", func(c *fiber.Ctx) error {
		return c.SendFile("./web/dist/index.html")
	})
}
