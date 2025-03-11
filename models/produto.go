package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enum para tipo de produto
type TipoProduto string

const (
	Fisico   TipoProduto = "fisico"
	Servico   TipoProduto = "servico"
)

// Produto é a estrutura base para todos os tipos de produtos
type Produto struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Nome        string         `json:"nome" validate:"required"`
	Descricao   string         `json:"descricao" validate:"required"`
	Preco       float64        `json:"preco" validate:"required,gt=0"`
	Tipo        TipoProduto    `json:"tipo" validate:"required,oneof=fisico servico"`
	Status      bool           `json:"status" gorm:"default:true"`
	CriadoEm    time.Time      `json:"criado_em"`
	AtualizadoEm time.Time     `json:"atualizado_em"`
}

// ProdutoFisico representa produtos físicos como roupas
type ProdutoFisico struct {
	Produto    Produto    `json:"produto" gorm:"embedded"`
	SKU         string    `json:"sku" validate:"required"`
	Estoque     int       `json:"estoque" validate:"required,gte=0"`
	Tamanho     string    `json:"tamanho,omitempty"`
	Cor         string    `json:"cor,omitempty"`
	Peso        float64   `json:"peso,omitempty" validate:"gte=0"`
}

// ProdutoServico representa serviços como assinaturas de academia
type ProdutoServico struct {
	Produto     Produto   `json:"produto" gorm:"embedded"`
	DuracaoMeses int       `json:"duracao_meses" validate:"required,gt=0"`
	Recorrente   bool      `json:"recorrente"`
	Beneficios   string    `json:"beneficios"`
}

// BeforeCreate será chamado antes de criar um novo produto
func (p *Produto) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	p.CriadoEm = time.Now()
	p.AtualizadoEm = time.Now()
	return nil
}

// BeforeUpdate será chamado antes de atualizar um produto
func (p *Produto) BeforeUpdate(tx *gorm.DB) error {
	p.AtualizadoEm = time.Now()
	return nil
}