package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/notification-service/internal/domain"
)

type PostgresNotificationRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresNotificationRepository(pool *pgxpool.Pool) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{pool: pool}
}

func (r *PostgresNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, notification_type, title, message, channel, priority, status, action_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query, notification.ID, notification.UserID, notification.NotificationType,
		notification.Title, notification.Message, notification.Channel, notification.Priority,
		notification.Status, notification.ActionURL, notification.CreatedAt)
	return err
}

func (r *PostgresNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	query := `SELECT id, user_id, notification_type, title, message, channel, status, read_at, created_at FROM notifications WHERE id = $1`

	n := &domain.Notification{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&n.ID, &n.UserID, &n.NotificationType, &n.Title, &n.Message, &n.Channel, &n.Status, &n.ReadAt, &n.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}
	return n, nil
}

func (r *PostgresNotificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Notification, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE user_id = $1", userID).Scan(&total)

	query := `SELECT id, user_id, notification_type, title, message, channel, status, created_at FROM notifications
	          WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notifications := []*domain.Notification{}
	for rows.Next() {
		n := &domain.Notification{}
		rows.Scan(&n.ID, &n.UserID, &n.NotificationType, &n.Title, &n.Message, &n.Channel, &n.Status, &n.CreatedAt)
		notifications = append(notifications, n)
	}
	return notifications, total, nil
}

func (r *PostgresNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE notifications SET status = 'read', read_at = $1 WHERE id = $2`
	result, err := r.pool.Exec(ctx, query, now, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

func (r *PostgresNotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND status != 'read'", userID).Scan(&count)
	return count, err
}

type PostgresMessageRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresMessageRepository(pool *pgxpool.Pool) *PostgresMessageRepository {
	return &PostgresMessageRepository{pool: pool}
}

func (r *PostgresMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `INSERT INTO messages (id, booking_id, sender_id, receiver_id, content, attachments, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, message.ID, message.BookingID, message.SenderID, message.ReceiverID,
		message.Content, message.Attachments, message.CreatedAt)
	return err
}

func (r *PostgresMessageRepository) GetByBooking(ctx context.Context, bookingID uuid.UUID, offset, limit int) ([]*domain.Message, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM messages WHERE booking_id = $1", bookingID).Scan(&total)

	query := `SELECT id, booking_id, sender_id, receiver_id, content, created_at FROM messages
	          WHERE booking_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, bookingID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	messages := []*domain.Message{}
	for rows.Next() {
		m := &domain.Message{}
		rows.Scan(&m.ID, &m.BookingID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt)
		messages = append(messages, m)
	}
	return messages, total, nil
}
