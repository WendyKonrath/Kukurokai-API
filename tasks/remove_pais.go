package tasks

import (
	"go-api/db"
	"go-api/models"
	"time"
)

func RemovePaisAutomaticamente() {
	for {
		// Encontrar clientes que completaram 18 anos
		var clientes []models.Cliente
		config.DB.Preload("Pais").Find(&clientes)

		for _, cliente := range clientes {
			idade := time.Now().Year() - cliente.DataNascimento.Year()
			if time.Now().Before(cliente.DataNascimento.AddDate(idade, 0, 0)) {
				idade-- // Ajuste para o caso de o aniversário ainda não ter ocorrido este ano
			}

			if idade >= 18 && cliente.PaisID != nil {
				// Remover os pais
				config.DB.Delete(&models.Pais{}, "cliente_id = ?", cliente.ID)
				cliente.PaisID = nil
				config.DB.Save(&cliente)
			}
		}

		// Executar a verificação diariamente
		time.Sleep(24 * time.Hour)
	}
}