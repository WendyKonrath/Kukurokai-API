package routes

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func SetupClienteRoutes(app *fiber.App) {
	clienteGroup := app.Group("/clientes", middleware.JWTMiddleware())

	clienteGroup.Get("/", GetClientes)
	clienteGroup.Post("/", CreateCliente)
}

func GetClientes(c *fiber.Ctx) error {
	var clientes []models.Cliente
	config.DB.Find(&clientes)
	return c.JSON(clientes)
}

func CreateCliente(c *fiber.Ctx) error {
	cliente := new(models.Cliente)

	if err := c.BodyParser(cliente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	if err := validate.Struct(cliente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	config.DB.Create(&cliente)
	return c.JSON(cliente)
}
