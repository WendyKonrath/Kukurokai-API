package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enum para método de pagamento
type PaymentMethod string

const (
	Boleto     PaymentMethod = "boleto"
	Pix        PaymentMethod = "pix"
	DebitCard  PaymentMethod = "debit_card"
	CreditCard PaymentMethod = "credit_card"
)

// Enum para status do pagamento
type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Paid      PaymentStatus = "paid"
	Overdue   PaymentStatus = "overdue"
	Cancelled PaymentStatus = "cancelled"
)

// Subscription representa uma assinatura de serviço
type Subscription struct {
	ID              string        `json:"id" gorm:"primaryKey"`
	ClienteID       string        `json:"cliente_id" validate:"required"`
	Cliente         Cliente       `json:"cliente" gorm:"foreignKey:ClienteID;references:ID"`
	PaymentMethod   PaymentMethod `json:"payment_method" validate:"required,oneof=boleto pix debit_card credit_card"`
	CardNumber      *string       `json:"card_number,omitempty" gorm:"type:text"` // Encrypted
	CardCVV         *string       `json:"card_cvv,omitempty" gorm:"type:text"`    // Encrypted
	BillingDay      int           `json:"billing_day" validate:"required,min=1,max=31"`
	PaymentStatus   PaymentStatus `json:"payment_status" gorm:"default:'pending'"`
	NextBillingDate time.Time     `json:"next_billing_date"`
	Amount          float64       `json:"amount" validate:"required,gt=0"`
	Active          bool          `json:"active" gorm:"default:true"`
	CriadoEm        time.Time     `json:"criado_em"`
	AtualizadoEm    time.Time     `json:"atualizado_em"`
}

// BeforeCreate será chamado antes de criar uma nova assinatura
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.CriadoEm = time.Now()
	s.AtualizadoEm = time.Now()

	// Configurar a próxima data de cobrança
	now := time.Now()
	s.NextBillingDate = time.Date(now.Year(), now.Month(), s.BillingDay, 0, 0, 0, 0, time.Local)
	if s.NextBillingDate.Before(now) {
		s.NextBillingDate = s.NextBillingDate.AddDate(0, 1, 0)
	}

	return nil
}

// BeforeUpdate será chamado antes de atualizar uma assinatura
func (s *Subscription) BeforeUpdate(tx *gorm.DB) error {
	s.AtualizadoEm = time.Now()
	return nil
}

