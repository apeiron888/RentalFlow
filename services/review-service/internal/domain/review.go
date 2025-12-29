package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReviewType string

const (
	TypeRenterToOwner ReviewType = "renter_to_owner"
	TypeOwnerToRenter ReviewType = "owner_to_renter"
	TypeRenterToItem  ReviewType = "renter_to_item"
)

type Review struct {
	ID           uuid.UUID  `json:"id"`
	BookingID    uuid.UUID  `json:"booking_id"`
	ReviewerID   uuid.UUID  `json:"reviewer_id"`
	TargetUserID *uuid.UUID `json:"target_user_id,omitempty"`
	TargetItemID *uuid.UUID `json:"target_item_id,omitempty"`
	ReviewType   ReviewType `json:"review_type"`
	Rating       float64    `json:"rating"`
	Comment      string     `json:"comment"`
	IsVerified   bool       `json:"is_verified"`
	IsVisible    bool       `json:"is_visible"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewReview(bookingID, reviewerID uuid.UUID, reviewType ReviewType, rating float64, comment string) *Review {
	if rating < 1.0 || rating > 5.0 {
		rating = 5.0
	}
	return &Review{
		ID:         uuid.New(),
		BookingID:  bookingID,
		ReviewerID: reviewerID,
		ReviewType: reviewType,
		Rating:     rating,
		Comment:    comment,
		IsVerified: false,
		IsVisible:  true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
