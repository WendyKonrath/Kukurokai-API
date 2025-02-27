package models

import (
	"time"

	"github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// Modelo de Cliente
type Cliente struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	Nome              string    `json:"nome" gorm:"type:varchar(255);not null"`
	DataNascimento    time.Time `json:"data_nascimento" gorm:"not null"`
	Genero            string    `json:"genero" gorm:"type:varchar(50)"`
	Email             string    `json:"email" gorm:"unique"`
	Telefone          string    `json:"telefone"`
	CPF               string    `json:"cpf" gorm:"unique"`
	Endereco          string    `json:"endereco"`
	Cidade            string    `json:"cidade"`
	Estado            string    `json:"estado"`
	CEP               string    `json:"cep"`
	FlagAniversariante bool     `json:"flag_aniversariante" gorm:"default:false"`
	FlagInadimplente   bool     `json:"flag_inadimplente" gorm:"default:false"`
	PaisID            *string   `json:"pais_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Modelo de Pais
type Pais struct {
	ID          string `json:"id" gorm:"primaryKey"`
	ClienteID   string `json:"cliente_id" gorm:"unique;constraint:OnDelete:CASCADE"`
	NomePai     string `json:"nome_pai"`
	TelefonePai string `json:"telefone_pai"`
	EmailPai    string `json:"email_pai"`
	CPFPai      string `json:"cpf_pai"`
	NomeMae     string `json:"nome_mae"`
	TelefoneMae string `json:"telefone_mae"`
	EmailMae    string `json:"email_mae"`
	CPFMae      string `json:"cpf_mae"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Gerar ID automaticamente com nanoid
func (c *Cliente) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID, err = gonanoid.New()
	}
	return
}

func (p *Pais) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID, err = gonanoid.New()
	}
	return
}

