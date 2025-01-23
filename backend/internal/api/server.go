package api

import (
	"log/slog"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	app    *fiber.App
	logger *slog.Logger
}

func NewServer(handler *Handler, logger *slog.Logger) *Server {
	s := &Server{
		logger: logger,
	}

	s.setupApp()

	handler.RegisterRoutes(s.app)

	return s
}

func (s *Server) setupApp() {
	s.app = fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	s.app.Use(recover.New())
	s.app.Use(requestid.New())
	s.app.Use(s.logMiddleware)
	s.app.Use(s.corsMiddleware())
}

func (s *Server) logMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	s.logger.Info("Request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Int("status", c.Response().StatusCode()),
		slog.Duration("duration", time.Since(start)),
	)
	return err
}

func (s *Server) corsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			if (u.Scheme == "http" || u.Scheme == "https") && (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1") {
				return true
			}
			return false
		},
		AllowHeaders: "Origin, Accept, Content-Type, nOrigin, X-Requested-With",
		AllowMethods: "GET,OPTIONS",
	})
}

func (s *Server) Start(port string) error {
	s.logger.Info("Starting server", "port", port)
	return s.app.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down server")
	return s.app.Shutdown()
}
