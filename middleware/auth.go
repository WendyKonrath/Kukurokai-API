package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken gera um token JWT com as claims necessárias
func GenerateToken(c *fiber.Ctx) error {
	claims := jwt.MapClaims{
		"user": "cliente_test",            // Nome ou identificação do usuário
		"role": "admin",                   // Adiciona uma role para controle de permissão
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Token expira em 72 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("sua_chave_secreta"))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao gerar token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

// JWTMiddleware valida o token JWT
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")

        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token não fornecido"})
        }

        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        // Verifica e valida o token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("sua_chave_secreta"), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido"})
        }

        // Acessa as claims e verifica a role
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token malformado"})
        }

        role := claims["role"].(string) // Supondo que a role esteja no token
        fmt.Println("Role do usuário:", role)

        // Verifique se o usuário tem permissão para acessar essa rota
        if role != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Acesso proibido"})
        }

        c.Locals("user", claims)

        return c.Next()
    }
}
