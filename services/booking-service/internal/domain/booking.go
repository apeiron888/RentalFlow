package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusActive    BookingStatus = "active"
	StatusCompleted BookingStatus = "completed"
	StatusCancelled BookingStatus = "cancelled"
)

type CancellationPolicy string

const (
	PolicyFlexible CancellationPolicy = "flexible"
	PolicyModerate CancellationPolicy = "moderate"
	PolicyStrict   CancellationPolicy = "strict"
)

type Booking struct {
	ID                 uuid.UUID
	BookingNumber      string
	RenterID           uuid.UUID
	OwnerID            uuid.UUID
	RentalItemID       uuid.UUID
	Status             BookingStatus
	StartDate          time.Time
	EndDate            time.Time
	TotalDays          int
	DailyRate          float64
	Subtotal           float64
	SecurityDeposit    float64
	ServiceFee         float64
	TotalAmount        float64
	PickupAddress      string
	PickupNotes        string
	PickupTime         *time.Time
	ReturnAddress      string
	ReturnNotes        string
	ReturnTime         *time.Time
	CancellationPolicy CancellationPolicy
	AgreementSigned    bool
	AgreementURL       string
	CancelledBy        *uuid.UUID
	CancellationReason string
	PaymentStatus      string
	PaymentID          *uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewBooking(renterID, ownerID, rentalItemID uuid.UUID, startDate, endDate time.Time, dailyRate, securityDeposit float64) *Booking {
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	if totalDays < 1 {
		totalDays = 1
	}

	subtotal := float64(totalDays) * dailyRate
	serviceFee := subtotal * 0.10 // 10%
	totalAmount := subtotal + serviceFee + securityDeposit

	now := time.Now()
	return &Booking{
		ID:                 uuid.New(),
		BookingNumber:      generateBookingNumber(),
		RenterID:           renterID,
		OwnerID:            ownerID,
		RentalItemID:       rentalItemID,
		Status:             StatusPending,
		StartDate:          startDate,
		EndDate:            endDate,
		TotalDays:          totalDays,
		DailyRate:          dailyRate,
		Subtotal:           subtotal,
		SecurityDeposit:    securityDeposit,
		ServiceFee:         serviceFee,
		TotalAmount:        totalAmount,
		CancellationPolicy: PolicyModerate,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func generateBookingNumber() string {
	return "BK" + time.Now().Format("20060102") + uuid.New().String()[:4]
}
