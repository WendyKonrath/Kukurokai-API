package models

import (
	"time"
	"github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Cliente struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	Nome              string    `json:"nome" validate:"required,min=3"`
	DataNascimento    time.Time `json:"data_nascimento" validate:"required"`
	Genero            string    `json:"genero" validate:"required,oneof=Masculino Feminino Outro"`
	Email             string    `json:"email" validate:"required,email"`
	Telefone          string    `json:"telefone" validate:"required"`
	CPF               string    `json:"cpf" validate:"required"`
	Endereco          string    `json:"endereco" validate:"required"`
	Cidade            string    `json:"cidade" validate:"required"`
	Estado            string    `json:"estado" validate:"required"`
	CEP               string    `json:"cep" validate:"required"`
	FlagAniversariante bool     `json:"flag_aniversariante"`
	FlagInadimplente   bool     `json:"flag_inadimplente"`
	PaisID            *string   `json:"pais_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Pais struct {
	ID          string `json:"id" gorm:"primaryKey"`
	ClienteID   string `json:"cliente_id" gorm:"unique;constraint:OnDelete:CASCADE"` // Chave estrangeira para Cliente
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
func (u *Pais) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID, err = gonanoid.New()
	}
	return
}

// Gerar ID automaticamente com nanoid
func (u *Cliente) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID, err = gonanoid.New()
	}
	return
}