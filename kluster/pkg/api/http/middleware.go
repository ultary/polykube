package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func NewAnalyzer(store *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		session, err := store.Get(ctx)
		if err != nil {
			return err
		}

		// TODO: trace user access before handle request

		if err = session.Save(); err != nil {
			log.Fatal(err)
		}
		return ctx.Next()
	}
}

const (
	hle = healthcheck.DefaultLivenessEndpoint
	hre = healthcheck.DefaultReadinessEndpoint
)

func NewLogger(isSkip func(ctx *fiber.Ctx) bool) fiber.Handler {
	l := logger.New()
	return func(ctx *fiber.Ctx) error {
		if isSkip(ctx) {
			return ctx.Next()
		}
		return l(ctx)
	}
}
