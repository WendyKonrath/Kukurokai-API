// routes/cliente.go
package routes

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var validate = validator.New()

func SetupClienteRoutes(app *fiber.App) {
	clienteGroup := app.Group("/clientes", middleware.JWTMiddleware())

	clienteGroup.Get("/", GetClientes)
	clienteGroup.Post("/", CreateCliente)
	clienteGroup.Put("/:id", UpdateCliente)
	clienteGroup.Delete("/:id", DeleteCliente)

	// Rota temporária para criar usuário
	app.Post("/create-user", CreateUser)
}

// Função para criar um usuário
func CreateUser(c *fiber.Ctx) error {
	type CreateUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Role     string `json:"role" validate:"required,oneof=superadmin admin user"`
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Criptografar a senha
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criptografar senha"})
	}

	// Criar o usuário
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar usuário"})
	}

	return c.JSON(fiber.Map{"message": "Usuário criado com sucesso", "user": user})
}

func GetClientes(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var clientes []models.Cliente
	config.DB.Preload("Pais").Find(&clientes) // Carrega os dados de Pais junto com Cliente
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

	// Validação dos dados do Pais
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