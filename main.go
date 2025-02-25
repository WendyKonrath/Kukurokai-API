package main

import (
	"log"

	"go-api/db"
	"go-api/middleware"
	"go-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitDB()

	app := fiber.New()
	middleware.SetupSecurity(app)
	routes.SetupClienteRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
