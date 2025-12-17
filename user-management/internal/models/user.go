package models

import (
	"context"
	"time"
)

type Basemodel struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt  time.Time `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

type UserStats struct {
	Orders        int       `json:"orders,omitempty" bson:"orders,omitempty"`
	Spent         int       `json:"spent,omitempty" bson:"spent,omitempty"`
	Reviews       int       `json:"reviews,omitempty" bson:"reviews,omitempty"`
	Since         time.Time `json:"since,omitempty" bson:"since,omitempty"`
	LastOrderDate time.Time `json:"last_order_date,omitempty" bson:"last_order_date,omitempty"`
}

type Themes string

const (
	Light Themes = "light"
	Dark  Themes = "Dark"
)

type UserPreferences struct {
	Language string `json:"language,omitempty" bson:"language,omitempty"`
	Currency string `json:"currency,omitempty" bson:"currency,omitempty"`
	TimeZone string `json:"time_zone,omitempty" bson:"time_zone,omitempty"`
	Theme    Themes `json:"theme,omitempty" bson:"theme,omitempty"`
}

type Address struct {
	Street  string `json:"street,omitempty" bson:"street,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	Country string `json:"country,omitempty" bson:"country,omitempty"`
}

type User struct {
	MetaData       Basemodel       `json:"meta_data,omitempty" bson:"meta_data,omitempty"`
	UserId         string          `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Email          string          `json:"email,omitempty" bson:"email,omitempty"`
	FullName       string          `json:"full_name,omitempty" bson:"full_name,omitempty"`
	PhoneNumber    string          `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	DOB            time.Time       `json:"dob,omitempty" bson:"dob,omitempty"`
	Gender         Sex             `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfilePicture string          `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
	Bio            string          `json:"bio,omitempty" bson:"bio,omitempty"`
	stats          UserStats       `json:"stats,omitempty" bson:"stats,omitempty"`
	preferences    UserPreferences `json:"preferences,omitempty" bson:"preferences,omitempty"`
	address        Address         `json:"address,omitempty" bson:"address,omitempty"`
}

type UserService interface {
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, id string, update User) (User, error)
	GetById(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	ListAll(ctx context.Context, filter User) (*User, error)
	Delete(ctx context.Context, id string) error
}

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, id string, update User) (User, error)
	GetById(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	ListAll(ctx context.Context, filter User) (*User, error)
	Delete(ctx context.Context, id string) error
}
