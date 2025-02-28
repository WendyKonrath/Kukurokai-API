// routes/cliente.go
package routes

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var validate = utils.Validate

func SetupClienteRoutes(app *fiber.App) {
	clienteGroup := app.Group("/clientes", middleware.JWTMiddleware())

	clienteGroup.Get("/", GetClientes)
	clienteGroup.Post("/", CreateCliente)
	clienteGroup.Put("/:id", UpdateCliente)
	clienteGroup.Delete("/:id", DeleteCliente)
}

func GetClientes(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var clientes []models.Cliente
	config.DB.Find(&clientes)
	return c.JSON(clientes)
}

func CreateCliente(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	type ClienteRequest struct {
		Cliente models.Cliente `json:"cliente"`
		Pais    models.Pais    `json:"pais"`
	}

	var req ClienteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados do cliente
	if err := validate.Struct(req.Cliente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Verificar se o cliente é menor de idade
	idade := time.Now().Year() - req.Cliente.DataNascimento.Year()
	if time.Now().Before(req.Cliente.DataNascimento.AddDate(idade, 0, 0)) {
		idade-- // Ajuste para o caso de o aniversário ainda não ter ocorrido este ano
	}

	if idade < 18 {
		// Validação dos dados do Pais (apenas se o cliente for menor de idade)
		if err := validate.Struct(req.Pais); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
		}

		// Ignorar o ID enviado pelo cliente e gerar um novo
		req.Cliente.ID = "" // Remove o ID enviado pelo cliente
		if err := config.DB.Create(&req.Cliente).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar cliente"})
		}

		// Associa o Pais ao Cliente
		req.Pais.ClienteID = req.Cliente.ID
		req.Pais.ID = "" // Remove o ID enviado pelo cliente
		if err := config.DB.Create(&req.Pais).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar pais"})
		}

		// Atualiza o PaisID no Cliente
		req.Cliente.PaisID = &req.Pais.ID
		if err := config.DB.Save(&req.Cliente).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar cliente"})
		}
	} else {
		// Cliente é maior de idade, não adiciona os pais
		req.Cliente.ID = "" // Remove o ID enviado pelo cliente
		if err := config.DB.Create(&req.Cliente).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar cliente"})
		}
	}

	return c.JSON(req.Cliente)
}

func UpdateCliente(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var cliente models.Cliente
	if err := config.DB.First(&cliente, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
	}

	if err := c.BodyParser(&cliente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados do cliente
	if err := validate.Struct(cliente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	config.DB.Save(&cliente)
	return c.JSON(cliente)
}

func DeleteCliente(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var cliente models.Cliente
	if err := config.DB.First(&cliente, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
	}

	config.DB.Delete(&cliente)
	return c.SendStatus(204)
}