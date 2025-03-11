package routes

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupSubscriptionRoutes(app *fiber.App) {
	subGroup := app.Group("/subscriptions", middleware.JWTMiddleware())

	// Rotas de assinatura
	subGroup.Get("/", ListSubscriptions)
	subGroup.Get("/:id", GetSubscription)
	subGroup.Post("/", CreateSubscription)
	subGroup.Put("/:id", UpdateSubscription)
	subGroup.Delete("/:id", CancelSubscription)
}

func ListSubscriptions(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var subscriptions []models.Subscription
	if err := config.DB.Find(&subscriptions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar assinaturas"})
	}

	return c.JSON(subscriptions)
}

func GetSubscription(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var subscription models.Subscription
	if err := config.DB.First(&subscription, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Assinatura não encontrada"})
	}

	return c.JSON(subscription)
}

func CreateSubscription(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var subscription models.Subscription
	if err := c.BodyParser(&subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados
	if err := utils.Validate.Struct(subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Criptografar dados do cartão se fornecidos
	if subscription.CardNumber != nil {
		encryptedNumber, err := utils.Encrypt([]byte(*subscription.CardNumber))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criptografar dados do cartão"})
		}
		subscription.CardNumber = &encryptedNumber
	}

	if subscription.CardCVV != nil {
		encryptedCVV, err := utils.Encrypt([]byte(*subscription.CardCVV))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criptografar dados do cartão"})
		}
		subscription.CardCVV = &encryptedCVV
	}

	// Configurar status inicial
	subscription.PaymentStatus = models.Pending

	// Criar assinatura
	if err := config.DB.Create(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar assinatura"})
	}

	return c.JSON(subscription)
}

func UpdateSubscription(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var subscription models.Subscription
	if err := config.DB.First(&subscription, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Assinatura não encontrada"})
	}

	if err := c.BodyParser(&subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados
	if err := utils.Validate.Struct(subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Atualizar assinatura
	if err := config.DB.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar assinatura"})
	}

	return c.JSON(subscription)
}

func CancelSubscription(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var subscription models.Subscription
	if err := config.DB.First(&subscription, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Assinatura não encontrada"})
	}

	// Marcar como cancelada e inativa
	subscription.PaymentStatus = models.Cancelled
	subscription.Active = false

	if err := config.DB.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao cancelar assinatura"})
	}

	return c.SendStatus(204)
}