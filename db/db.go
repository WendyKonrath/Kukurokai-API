package config

import (
	"fmt"
	"go-api/db/config"
	"go-api/models" // Certifique-se de importar o pacote onde estão seus modelos
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	config.LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Adiciona a migração aqui
	if err := DB.AutoMigrate(&models.Cliente{}, &models.Pais{}, &models.User{}, &models.Sale{}, &models.Produto{}, &models.Subscription{}); err != nil {
		log.Fatal("Erro ao migrar as tabelas:", err)
	}

	fmt.Println("Banco de dados conectado e migrações executadas com sucesso!")
}
