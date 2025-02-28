// main.go
package main

import (
	"go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/routes"
	"go-api/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitDB()

	// Verificar e criar o superadmin
	createSuperAdmin()

	app := fiber.New()
	middleware.SetupSecurity(app)
	routes.SetupAuthRoutes(app)
	routes.SetupClienteRoutes(app)
	routes.SetupUserRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

func createSuperAdmin() {
	superAdminEmail := os.Getenv("SUPERADMIN_EMAIL")
	superAdminPassword := os.Getenv("SUPERADMIN_PASSWORD")

	if superAdminEmail == "" || superAdminPassword == "" {
		log.Fatal("Variáveis de ambiente SUPERADMIN_EMAIL e SUPERADMIN_PASSWORD não configuradas")
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", superAdminEmail).First(&existingUser).Error; err == nil {
		log.Println("Superadmin já existe no banco de dados")
		return
	}

	hashedPassword, err := utils.HashPassword(superAdminPassword)
	if err != nil {
		log.Fatal("Erro ao criptografar senha do superadmin")
	}

	superAdmin := models.User{
		Email:    superAdminEmail,
		Password: hashedPassword,
		Role:     "superadmin",
	}

	if err := config.DB.Create(&superAdmin).Error; err != nil {
		log.Fatal("Erro ao criar superadmin no banco de dados")
	}

	log.Println("Superadmin criado com sucesso")
}