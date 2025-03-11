package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sale representa uma venda de produto
type Sale struct {
	ID              string         `json:"id" gorm:"primaryKey"`
	ProdutoID       string         `json:"produto_id" validate:"required"`
	Produto         Produto        `json:"produto" gorm:"foreignKey:ProdutoID;references:ID"`
	Valor           float64        `json:"valor" validate:"required,gt=0"`
	Custo           float64        `json:"custo" validate:"required,gte=0"`
	LucroLiquido    float64        `json:"lucro_liquido"`
	Pago            bool           `json:"pago" gorm:"default:false"`
	ClienteID       string         `json:"cliente_id" validate:"required"`
	Cliente         Cliente        `json:"cliente" gorm:"foreignKey:ClienteID;references:ID"`
	Quantidade      int            `json:"quantidade" validate:"required,gt=0"`
	FormaPagamento  PaymentMethod  `json:"forma_pagamento" validate:"required,oneof=boleto pix debit_card credit_card"`
	Status          PaymentStatus  `json:"status" gorm:"default:'pending'"`
	CriadoEm        time.Time      `json:"criado_em"`
	AtualizadoEm    time.Time     `json:"atualizado_em"`
}

// BeforeCreate será chamado antes de criar uma nova venda
func (s *Sale) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	
	// Calcular o lucro líquido
	s.LucroLiquido = s.Valor - s.Custo
	
	// Se estiver marcado como pago, atualizar o status
	if s.Pago {
		s.Status = Paid
	}
	
	s.CriadoEm = time.Now()
	s.AtualizadoEm = time.Now()
	return nil
}

// BeforeUpdate será chamado antes de atualizar uma venda
func (s *Sale) BeforeUpdate(tx *gorm.DB) error {
	// Recalcular o lucro líquido em caso de alteração nos valores
	s.LucroLiquido = s.Valor - s.Custo
	
	// Atualizar o status de pagamento
	if s.Pago && s.Status != Cancelled {
		s.Status = Paid
	}
	
	s.AtualizadoEm = time.Now()
	return nil
}