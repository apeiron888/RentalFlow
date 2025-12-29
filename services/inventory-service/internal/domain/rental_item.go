package domain

import (
	"time"

	"github.com/google/uuid"
)

// ItemCategory represents the category of a rental item
type ItemCategory string

const (
	CategoryVehicle   ItemCategory = "vehicle"
	CategoryEquipment ItemCategory = "equipment"
	CategoryProperty  ItemCategory = "property"
)

// IsValid checks if the category is valid
func (c ItemCategory) IsValid() bool {
	switch c {
	case CategoryVehicle, CategoryEquipment, CategoryProperty:
		return true
	}
	return false
}

// AvailabilityStatus represents the availability status
type AvailabilityStatus string

const (
	StatusAvailable   AvailabilityStatus = "available"
	StatusBooked      AvailabilityStatus = "booked"
	StatusMaintenance AvailabilityStatus = "maintenance"
	StatusBlocked     AvailabilityStatus = "blocked"
)

// MaintenanceStatus represents maintenance status
type MaintenanceStatus string

const (
	MaintenanceScheduled  MaintenanceStatus = "scheduled"
	MaintenanceInProgress MaintenanceStatus = "in_progress"
	MaintenanceCompleted  MaintenanceStatus = "completed"
)

// RentalItem represents a rental item
type RentalItem struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Title       string
	Description string
	Category    ItemCategory
	Subcategory string

	// Pricing
	DailyRate       float64
	WeeklyRate      float64
	MonthlyRate     float64
	SecurityDeposit float64

	// Location
	Address   string
	City      string
	Latitude  float64
	Longitude float64

	// Specifications (stored as map)
	Specifications map[string]string

	// Images
	Images []string

	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewRentalItem creates a new rental item
func NewRentalItem(ownerID uuid.UUID, title, description string, category ItemCategory, subcategory string) *RentalItem {
	now := time.Now()
	return &RentalItem{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		Title:          title,
		Description:    description,
		Category:       category,
		Subcategory:    subcategory,
		Specifications: make(map[string]string),
		Images:         []string{},
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// AvailabilitySlot represents an availability slot for a rental item
type AvailabilitySlot struct {
	ID           uuid.UUID
	RentalItemID uuid.UUID
	StartDate    time.Time
	EndDate      time.Time
	Status       AvailabilityStatus
	BookingID    *uuid.UUID
	CreatedAt    time.Time
}

// NewAvailabilitySlot creates a new availability slot
func NewAvailabilitySlot(rentalItemID uuid.UUID, startDate, endDate time.Time, status AvailabilityStatus) *AvailabilitySlot {
	return &AvailabilitySlot{
		ID:           uuid.New(),
		RentalItemID: rentalItemID,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		CreatedAt:    time.Now(),
	}
}

// MaintenanceLog represents a maintenance log entry
type MaintenanceLog struct {
	ID              uuid.UUID
	RentalItemID    uuid.UUID
	MaintenanceType string
	Description     string
	StartDate       time.Time
	EndDate         *time.Time
	Cost            float64
	Status          MaintenanceStatus
	CreatedAt       time.Time
}

// NewMaintenanceLog creates a new maintenance log
func NewMaintenanceLog(rentalItemID uuid.UUID, maintenanceType, description string, startDate time.Time, cost float64) *MaintenanceLog {
	return &MaintenanceLog{
		ID:              uuid.New(),
		RentalItemID:    rentalItemID,
		MaintenanceType: maintenanceType,
		Description:     description,
		StartDate:       startDate,
		Cost:            cost,
		Status:          MaintenanceScheduled,
		CreatedAt:       time.Now(),
	}
}

// Location represents a geographic location
type Location struct {
	Address   string
	City      string
	Latitude  float64
	Longitude float64
}
