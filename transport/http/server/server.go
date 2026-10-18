package server

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/noonbyte/platform/api"
	"github.com/noonbyte/platform/configs"
	"github.com/noonbyte/platform/middlewares"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewApp(c configs.AppConfiguration) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit:    20 * 1024 * 1024,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,

		TrustProxy: true,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: false,
		PanicHandler: func(c fiber.Ctx, recovered any) error {
			var err error
			stack := debug.Stack()

			event := log.Error().
				Str("method", c.Method()).
				Str("path", c.Path()).
				Str("ip", c.IP()).
				Str("request_id", c.Get("X-Request-ID"))

			err, ok := recovered.(error)
			if ok {
				event = event.Err(err)
			} else {
				event = event.Interface("panic", recovered)
			}

			event.Msgf(
				"Unhandled panic: %v\n\nStack trace:\n%s",
				recovered,
				stack,
			)

			return api.InternalServerError(c, err)
		},
	}))

	app.Use(middlewares.Logging())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": fiber.Map{
				"status":    "ok",
				"timestamp": time.Now().Unix(),
			},
		})
	})

	return app
}

func Listen(ctx context.Context, app *fiber.App, c *configs.AppConfiguration, log *zerolog.Logger) error {

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	log.Printf("Server starting on http://%s", addr)

	go func() {
		<-ctx.Done()
		log.Print("Shutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Err(fmt.Errorf("Error during shutdown: %v", err))
		} else {
			log.Print("Server stopped gracefully")
		}
	}()

	if err := app.Listen(addr, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		log.Err(fmt.Errorf("Failed to start server: %v", err))
		return err
	}

	return nil
}
