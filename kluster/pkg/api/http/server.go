package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	app *fiber.App
}

func NewServer() *Server {

	app := fiber.New()

	app.Use(recover.New())
	app.Use(requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		Generator:  utils.UUIDv4,
		ContextKey: ContextKeyRequestID,
	}))
	app.Use(NewLogger(func(ctx *fiber.Ctx) bool {
		path := ctx.Path()
		return path == PathMetrics ||
			path == healthcheck.DefaultLivenessEndpoint ||
			path == healthcheck.DefaultReadinessEndpoint
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(healthcheck.New())
	app.Use(pprof.New())
	app.Get(PathMetrics, monitor.New())

	return &Server{
		app: app,
	}
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	log.Info("[HTTP] Stopping server")

	const timeout = 10 * time.Second
	err := s.app.ShutdownWithTimeout(timeout)

	log.Info("[HTTP] Stopped server")
	return err
}
