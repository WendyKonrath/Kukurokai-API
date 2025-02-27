package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// security.go
func SetupSecurity(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New())
	// app.Use(csrf.New()) // Desative temporariamente para testar
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60 * time.Second,
	}))
}

