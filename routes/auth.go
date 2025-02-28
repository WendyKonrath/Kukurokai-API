// routes/auth.go
package routes

import (
	"fmt"
	config "go-api/db"
	"go-api/models"
	"go-api/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/login", Login)
}

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciais inválidas"})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciais inválidas"})
	}

	// Verificar os dados antes de criptografar
	fmt.Println("Email:", user.Email)
	fmt.Println("Role:", user.Role)

	// Criptografar os dados sensíveis
	encryptedEmail, err := utils.Encrypt([]byte(user.Email))
	if err != nil {
		fmt.Println("Erro ao criptografar email:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criptografar dados"})
	}

	encryptedRole, err := utils.Encrypt([]byte(user.Role))
	if err != nil {
		fmt.Println("Erro ao criptografar role:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criptografar dados"})
	}

	// Criar o token com os dados criptografados
	claims := jwt.MapClaims{
		"user": encryptedEmail,
		"role": encryptedRole,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Chave secreta JWT não configurada"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao gerar token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}