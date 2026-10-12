package server

import (
	"context"
	"errors"
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
	})

	app.Use(recover.New(recover.Config{
		PanicHandler: func(c fiber.Ctx, r any) error {
			stack := debug.Stack()

			if err, ok := r.(error); ok {
				log.Error().
					Err(err).
					Str("panic", "recovered").
					Str("method", c.Method()).
					Str("path", c.Path()).
					Bytes("stack", stack).
					Msg("Unhandled panic")

				return api.InternalServerError(c, err)
			}
			log.Error().
				Interface("panic", r).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Bytes("stack", stack).
				Msg("Unhandled panic")

			return api.InternalServerError(c, errors.New("unknown error"))
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
