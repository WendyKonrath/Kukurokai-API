package routes

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupSaleRoutes(app *fiber.App) {
	saleGroup := app.Group("/sales", middleware.JWTMiddleware())

	// Rotas de venda
	saleGroup.Get("/", ListSales)
	saleGroup.Get("/:id", GetSale)
	saleGroup.Post("/", CreateSale)
	saleGroup.Put("/:id", UpdateSale)
	saleGroup.Delete("/:id", DeleteSale)
}

func ListSales(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var sales []models.Sale
	if err := config.DB.Preload("Produto").Preload("Cliente").Find(&sales).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar vendas"})
	}

	return c.JSON(sales)
}

func GetSale(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var sale models.Sale
	if err := config.DB.Preload("Produto").Preload("Cliente").First(&sale, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Venda não encontrada"})
	}

	return c.JSON(sale)
}

func CreateSale(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var sale models.Sale
	if err := c.BodyParser(&sale); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados
	if err := utils.Validate.Struct(sale); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Criar a venda
	if err := config.DB.Create(&sale).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar venda"})
	}

	return c.Status(201).JSON(sale)
}

func UpdateSale(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var sale models.Sale
	if err := config.DB.First(&sale, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Venda não encontrada"})
	}

	// Parse do corpo da requisição
	if err := c.BodyParser(&sale); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados
	if err := utils.Validate.Struct(sale); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Atualizar a venda
	if err := config.DB.Save(&sale).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar venda"})
	}

	return c.JSON(sale)
}

func DeleteSale(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var sale models.Sale
	if err := config.DB.First(&sale, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Venda não encontrada"})
	}

	// Deletar a venda
	if err := config.DB.Delete(&sale).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao deletar venda"})
	}

	return c.Status(204).Send(nil)
}