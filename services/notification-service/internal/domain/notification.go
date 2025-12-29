package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
	StatusRead    NotificationStatus = "read"
)

type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelSMS   NotificationChannel = "sms"
	ChannelPush  NotificationChannel = "push"
	ChannelInApp NotificationChannel = "in_app"
)

type Notification struct {
	ID               uuid.UUID           `json:"id"`
	UserID           uuid.UUID           `json:"user_id"`
	NotificationType string              `json:"notification_type"`
	Title            string              `json:"title"`
	Message          string              `json:"message"`
	Channel          NotificationChannel `json:"channel"`
	Priority         string              `json:"priority"`
	Status           NotificationStatus  `json:"status"`
	ActionURL        string              `json:"action_url"`
	SentAt           *time.Time          `json:"sent_at"`
	ReadAt           *time.Time          `json:"read_at"`
	CreatedAt        time.Time           `json:"created_at"`
}

type Message struct {
	ID          uuid.UUID  `json:"id"`
	BookingID   uuid.UUID  `json:"booking_id"`
	SenderID    uuid.UUID  `json:"sender_id"`
	ReceiverID  uuid.UUID  `json:"receiver_id"`
	Content     string     `json:"content"`
	Attachments []string   `json:"attachments"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func NewNotification(userID uuid.UUID, notifType, title, message string, channel NotificationChannel) *Notification {
	return &Notification{
		ID:               uuid.New(),
		UserID:           userID,
		NotificationType: notifType,
		Title:            title,
		Message:          message,
		Channel:          channel,
		Priority:         "medium",
		Status:           StatusPending,
		CreatedAt:        time.Now(),
	}
}

func NewMessage(bookingID, senderID, receiverID uuid.UUID, content string) *Message {
	return &Message{
		ID:         uuid.New(),
		BookingID:  bookingID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
}
