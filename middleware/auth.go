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
		"user": "cliente_test",            
		"role": "admin",                   
		
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("sua_chave_secreta"))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao gerar token"})
	}

	// Pegar o tempo de expiração direto das claims
	if exp, ok := claims["exp"].(int64); ok {
		expTime := time.Unix(exp, 0) // Converte o timestamp para time.Time
		fmt.Println("Token expira em:", expTime.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("Erro ao pegar tempo de expiração")
	}

	return c.JSON(fiber.Map{"token": tokenString})
}



// JWTMiddleware valida o token JWT
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Token não fornecido"})
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("sua_chave_secreta"), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Token malformado"})
		}

		c.Locals("user", claims)
		return c.Next()
	}
}