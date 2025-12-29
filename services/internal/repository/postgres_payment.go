package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/payment-service/internal/domain"
)

type PostgresPaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPaymentRepository(pool *pgxpool.Pool) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{pool: pool}
}

func (r *PostgresPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (
			id, booking_id, user_id, payment_type, amount, currency, status, method,
			rental_fee, security_deposit, service_fee, additional_services, tax,
			provider_name, checkout_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.pool.Exec(ctx, query,
		payment.ID, payment.BookingID, payment.UserID, payment.PaymentType, payment.Amount,
		payment.Currency, payment.Status, payment.Method,
		payment.RentalFee, payment.SecurityDeposit, payment.ServiceFee,
		payment.AdditionalServices, payment.Tax,
		payment.ProviderName, payment.CheckoutURL,
		payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

func (r *PostgresPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	query := `
		SELECT id, booking_id, user_id, payment_type, amount, currency, status, method,
			   rental_fee, security_deposit, service_fee, provider_transaction_id,
			   created_at, updated_at
		FROM payments WHERE id = $1
	`

	payment := &domain.Payment{}
	var providerTxID *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&payment.ID, &payment.BookingID, &payment.UserID, &payment.PaymentType,
		&payment.Amount, &payment.Currency, &payment.Status, &payment.Method,
		&payment.RentalFee, &payment.SecurityDeposit, &payment.ServiceFee,
		&providerTxID, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	if providerTxID != nil {
		payment.ProviderTransactionID = *providerTxID
	}

	return payment, nil
}

func (r *PostgresPaymentRepository) GetByBooking(ctx context.Context, bookingID uuid.UUID) ([]*domain.Payment, error) {
	query := `
		SELECT id, booking_id, user_id, payment_type, amount, currency, status, method,
			   created_at, updated_at
		FROM payments WHERE booking_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []*domain.Payment{}
	for rows.Next() {
		p := &domain.Payment{}
		err := rows.Scan(
			&p.ID, &p.BookingID, &p.UserID, &p.PaymentType, &p.Amount,
			&p.Currency, &p.Status, &p.Method, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PostgresPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments SET
			status = $2, provider_transaction_id = $3, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, payment.ID, payment.Status, payment.ProviderTransactionID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}
	return nil
}
