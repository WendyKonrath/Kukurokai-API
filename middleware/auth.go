// middleware/auth.go
package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go-api/utils"
)

// Carregar variáveis de ambiente
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Aviso: Arquivo .env não encontrado, usando variáveis do sistema")
	}
}

// JWTMiddleware valida o token JWT
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obter o token do header Authorization
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token não fornecido"})
		}

		// Remover o prefixo "Bearer " se presente
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validar o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar o algoritmo de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
			}

			// Obter a chave secreta do ambiente
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				return nil, fmt.Errorf("chave secreta JWT não configurada")
			}

			return []byte(jwtSecret), nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Token inválido", "details": err.Error()})
		}

		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
		}

		// Verificar as claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token malformado"})
		}

		// Descriptografar os dados sensíveis
		encryptedEmail := claims["user"].(string)
		email, err := utils.Decrypt(encryptedEmail)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Erro ao descriptografar dados"})
		}

		encryptedRole := claims["role"].(string)
		role, err := utils.Decrypt(encryptedRole)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Erro ao descriptografar dados"})
		}

		// Adicionar as claims descriptografadas ao contexto
		c.Locals("user", jwt.MapClaims{
			"email": string(email),
			"role":  string(role),
		})

		return c.Next()
	}
}