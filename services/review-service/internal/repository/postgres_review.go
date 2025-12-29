package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/review-service/internal/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.Review, int, error)
	GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Review, int, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostgresReviewRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresReviewRepository(pool *pgxpool.Pool) *PostgresReviewRepository {
	return &PostgresReviewRepository{pool: pool}
}

func (r *PostgresReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	query := `INSERT INTO reviews (id, booking_id, reviewer_id, target_user_id, target_item_id, review_type, rating, comment, is_verified, is_visible, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.pool.Exec(ctx, query, review.ID, review.BookingID, review.ReviewerID, review.TargetUserID, review.TargetItemID,
		review.ReviewType, review.Rating, review.Comment, review.IsVerified, review.IsVisible, review.CreatedAt, review.UpdatedAt)
	return err
}

func (r *PostgresReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	query := `SELECT id, booking_id, reviewer_id, target_user_id, target_item_id, review_type, rating, comment, is_visible, created_at, updated_at
	          FROM reviews WHERE id = $1`
	review := &domain.Review{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&review.ID, &review.BookingID, &review.ReviewerID, &review.TargetUserID,
		&review.TargetItemID, &review.ReviewType, &review.Rating, &review.Comment, &review.IsVisible, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReviewNotFound
		}
		return nil, err
	}
	return review, nil
}

func (r *PostgresReviewRepository) GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.Review, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM reviews WHERE target_item_id = $1 AND is_visible = true", itemID).Scan(&total)

	query := `SELECT id, booking_id, reviewer_id, review_type, rating, comment, created_at FROM reviews
	          WHERE target_item_id = $1 AND is_visible = true ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, itemID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reviews := []*domain.Review{}
	for rows.Next() {
		rev := &domain.Review{}
		rows.Scan(&rev.ID, &rev.BookingID, &rev.ReviewerID, &rev.ReviewType, &rev.Rating, &rev.Comment, &rev.CreatedAt)
		reviews = append(reviews, rev)
	}
	return reviews, total, nil
}

func (r *PostgresReviewRepository) GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Review, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM reviews WHERE target_user_id = $1 AND is_visible = true", userID).Scan(&total)

	query := `SELECT id, booking_id, reviewer_id, review_type, rating, comment, created_at FROM reviews
	          WHERE target_user_id = $1 AND is_visible = true ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reviews := []*domain.Review{}
	for rows.Next() {
		rev := &domain.Review{}
		rows.Scan(&rev.ID, &rev.BookingID, &rev.ReviewerID, &rev.ReviewType, &rev.Rating, &rev.Comment, &rev.CreatedAt)
		reviews = append(reviews, rev)
	}
	return reviews, total, nil
}

func (r *PostgresReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	query := `UPDATE reviews SET rating = $2, comment = $3, is_visible = $4, updated_at = NOW() WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, review.ID, review.Rating, review.Comment, review.IsVisible)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrReviewNotFound
	}
	return nil
}

func (r *PostgresReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM reviews WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrReviewNotFound
	}
	return nil
}
