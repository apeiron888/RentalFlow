package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/booking-service/internal/domain"
)

type PostgresBookingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresBookingRepository(pool *pgxpool.Pool) *PostgresBookingRepository {
	return &PostgresBookingRepository{pool: pool}
}

func (r *PostgresBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (
			id, booking_number, renter_id, owner_id, rental_item_id, status,
			start_date, end_date, total_days, daily_rate, subtotal,
			security_deposit, service_fee, total_amount, cancellation_policy,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.pool.Exec(ctx, query,
		booking.ID, booking.BookingNumber, booking.RenterID, booking.OwnerID, booking.RentalItemID, booking.Status,
		booking.StartDate, booking.EndDate, booking.TotalDays, booking.DailyRate, booking.Subtotal,
		booking.SecurityDeposit, booking.ServiceFee, booking.TotalAmount, booking.CancellationPolicy,
		booking.CreatedAt, booking.UpdatedAt,
	)
	return err
}

func (r *PostgresBookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	query := `
		SELECT id, booking_number, renter_id, owner_id, rental_item_id, status,
			   start_date, end_date, total_days, daily_rate, subtotal,
			   security_deposit, service_fee, total_amount, cancellation_policy,
			   agreement_signed, created_at, updated_at
		FROM bookings WHERE id = $1
	`

	booking := &domain.Booking{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&booking.ID, &booking.BookingNumber, &booking.RenterID, &booking.OwnerID, &booking.RentalItemID, &booking.Status,
		&booking.StartDate, &booking.EndDate, &booking.TotalDays, &booking.DailyRate, &booking.Subtotal,
		&booking.SecurityDeposit, &booking.ServiceFee, &booking.TotalAmount, &booking.CancellationPolicy,
		&booking.AgreementSigned, &booking.CreatedAt, &booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return booking, nil
}

func (r *PostgresBookingRepository) GetByRenter(ctx context.Context, renterID uuid.UUID, offset, limit int) ([]*domain.Booking, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings WHERE renter_id = $1", renterID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, booking_number, renter_id, owner_id, rental_item_id, status,
			   start_date, end_date, total_days, daily_rate, total_amount,
			   created_at, updated_at
		FROM bookings WHERE renter_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, renterID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	bookings := []*domain.Booking{}
	for rows.Next() {
		b := &domain.Booking{}
		err := rows.Scan(&b.ID, &b.BookingNumber, &b.RenterID, &b.OwnerID, &b.RentalItemID, &b.Status,
			&b.StartDate, &b.EndDate, &b.TotalDays, &b.DailyRate, &b.TotalAmount, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, b)
	}

	return bookings, total, nil
}

func (r *PostgresBookingRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]*domain.Booking, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings WHERE owner_id = $1", ownerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, booking_number, renter_id, owner_id, rental_item_id, status,
			   start_date, end_date, total_days, daily_rate, total_amount,
			   created_at, updated_at
		FROM bookings WHERE owner_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	bookings := []*domain.Booking{}
	for rows.Next() {
		b := &domain.Booking{}
		err := rows.Scan(&b.ID, &b.BookingNumber, &b.RenterID, &b.OwnerID, &b.RentalItemID, &b.Status,
			&b.StartDate, &b.EndDate, &b.TotalDays, &b.DailyRate, &b.TotalAmount, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, b)
	}

	return bookings, total, nil
}

func (r *PostgresBookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	query := `
		UPDATE bookings SET
			status = $2, agreement_signed = $3, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, booking.ID, booking.Status, booking.AgreementSigned)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrBookingNotFound
	}
	return nil
}
