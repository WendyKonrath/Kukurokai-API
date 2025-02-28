package routes

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/users", middleware.JWTMiddleware())

	userGroup.Get("/", ListUsers)       // Nova rota para listar usuários
	userGroup.Post("/", CreateUser)
	userGroup.Put("/:id", UpdateUser)
	userGroup.Delete("/:id", DeleteUser)
}

// Função para listar todos os usuários
func ListUsers(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	// Apenas o superadmin pode listar usuários
	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar usuários"})
	}

	return c.JSON(users)
}

// Função para criar um usuário
func CreateUser(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

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
	newUser := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar usuário"})
	}

	return c.JSON(fiber.Map{"message": "Usuário criado com sucesso", "user": newUser})
}

// Função para atualizar um usuário
func UpdateUser(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var existingUser models.User
	if err := config.DB.First(&existingUser, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuário não encontrado"})
	}

	type UpdateUserRequest struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"omitempty,min=6"`
		Role     string `json:"role" validate:"omitempty,oneof=superadmin admin user"`
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Atualizar os campos fornecidos
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criptografar senha"})
		}
		existingUser.Password = hashedPassword
	}
	if req.Role != "" {
		existingUser.Role = req.Role
	}

	if err := config.DB.Save(&existingUser).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar usuário"})
	}

	return c.JSON(fiber.Map{"message": "Usuário atualizado com sucesso", "user": existingUser})
}

// Função para deletar um usuário
func DeleteUser(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var existingUser models.User
	if err := config.DB.First(&existingUser, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuário não encontrado"})
	}

	if err := config.DB.Delete(&existingUser).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao deletar usuário"})
	}

	return c.SendStatus(204)
}