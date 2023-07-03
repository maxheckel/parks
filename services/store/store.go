package store

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

var Store = session.New(session.Config{
	Expiration: time.Hour * 24 * 365,
})
