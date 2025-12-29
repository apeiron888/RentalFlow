package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleRenter UserRole = "renter"
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
)

// IsValid checks if the role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleRenter, RoleOwner, RoleAdmin:
		return true
	}
	return false
}

// VerificationStatus represents the identity verification status
type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationVerified VerificationStatus = "verified"
	VerificationRejected VerificationStatus = "rejected"
)

// User represents a user in the system
type User struct {
	ID                    uuid.UUID
	Email                 string
	PasswordHash          string
	FirstName             string
	LastName              string
	Phone                 string
	Bio                   string
	AvatarURL             string
	Role                  UserRole
	IdentityVerified      bool
	VerificationStatus    VerificationStatus
	RefreshTokenHash      string
	RefreshTokenExpiresAt *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// NewUser creates a new user with default values
func NewUser(email, passwordHash, firstName, lastName, phone string, role UserRole) *User {
	now := time.Now()
	return &User{
		ID:                 uuid.New(),
		Email:              email,
		PasswordHash:       passwordHash,
		FirstName:          firstName,
		LastName:           lastName,
		Phone:              phone,
		Role:               role,
		IdentityVerified:   false,
		VerificationStatus: VerificationPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsOwner checks if the user is an owner
func (u *User) IsOwner() bool {
	return u.Role == RoleOwner
}

// IsRenter checks if the user is a renter
func (u *User) IsRenter() bool {
	return u.Role == RoleRenter
}

// CanManageItem checks if the user can manage rental items
func (u *User) CanManageItem() bool {
	return u.Role == RoleOwner || u.Role == RoleAdmin
}

// HasValidRefreshToken checks if the user has a valid refresh token
func (u *User) HasValidRefreshToken() bool {
	if u.RefreshTokenHash == "" || u.RefreshTokenExpiresAt == nil {
		return false
	}
	return u.RefreshTokenExpiresAt.After(time.Now())
}

// SetRefreshToken sets a new refresh token hash and expiry
func (u *User) SetRefreshToken(hash string, expiresAt time.Time) {
	u.RefreshTokenHash = hash
	u.RefreshTokenExpiresAt = &expiresAt
	u.UpdatedAt = time.Now()
}

// ClearRefreshToken clears the refresh token (logout)
func (u *User) ClearRefreshToken() {
	u.RefreshTokenHash = ""
	u.RefreshTokenExpiresAt = nil
	u.UpdatedAt = time.Now()
}

// IdentityDocument represents an identity document
type IdentityDocument struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	DocumentType string
	DocumentURL  string
	UploadedAt   time.Time
}

// DocumentType constants
const (
	DocTypeDriverLicense = "driver_license"
	DocTypeNationalID    = "national_id"
	DocTypePassport      = "passport"
)

// NewIdentityDocument creates a new identity document
func NewIdentityDocument(userID uuid.UUID, docType, docURL string) *IdentityDocument {
	return &IdentityDocument{
		ID:           uuid.New(),
		UserID:       userID,
		DocumentType: docType,
		DocumentURL:  docURL,
		UploadedAt:   time.Now(),
	}
}
