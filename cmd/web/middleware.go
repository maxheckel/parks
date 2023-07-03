package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxheckel/parks/services/store"
	"net/http"
)

func Authenticated(ctx *fiber.Ctx) error {
	session, err := store.Store.Get(ctx)
	if err != nil {
		return err
	}
	if session.Get("profile") == nil {
		return ctx.Redirect("/", http.StatusSeeOther)
	}
	return ctx.Next()
}
