package main

import (
	"fmt"
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
	app.Get("/generate-token", middleware.GenerateToken)
	app.Get("/test", middleware.JWTMiddleware(), func(c *fiber.Ctx) error {
		fmt.Println("Rota '/test' acessada com sucesso!")
		return c.SendString("Rota de teste acessada com sucesso!")
	})
	
	routes.SetupClienteRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
