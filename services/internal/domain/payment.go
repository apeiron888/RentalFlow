package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusCompleted  PaymentStatus = "completed"
	StatusFailed     PaymentStatus = "failed"
	StatusRefunded   PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	MethodChapa        PaymentMethod = "chapa"
	MethodTelebirr     PaymentMethod = "telebirr"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodCash         PaymentMethod = "cash"
)

type Payment struct {
	ID                    uuid.UUID     `json:"id"`
	BookingID             uuid.UUID     `json:"booking_id"`
	UserID                uuid.UUID     `json:"user_id"`
	PaymentType           string        `json:"payment_type"`
	Amount                float64       `json:"amount"`
	Currency              string        `json:"currency"`
	Status                PaymentStatus `json:"status"`
	Method                PaymentMethod `json:"method"`
	RentalFee             float64       `json:"rental_fee"`
	SecurityDeposit       float64       `json:"security_deposit"`
	ServiceFee            float64       `json:"service_fee"`
	AdditionalServices    float64       `json:"additional_services"`
	Tax                   float64       `json:"tax"`
	DepositHeld           bool          `json:"deposit_held"`
	DepositStatus         string        `json:"deposit_status"`
	ProviderName          string        `json:"provider_name"`
	ProviderTransactionID string        `json:"provider_transaction_id"`
	CheckoutURL           string        `json:"checkout_url"`
	ReceiptURL            string        `json:"receipt_url"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

func NewPayment(bookingID, userID uuid.UUID, amount float64, method PaymentMethod) *Payment {
	return &Payment{
		ID:        uuid.New(),
		BookingID: bookingID,
		UserID:    userID,
		Amount:    amount,
		Currency:  "ETB",
		Status:    StatusPending,
		Method:    method,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
