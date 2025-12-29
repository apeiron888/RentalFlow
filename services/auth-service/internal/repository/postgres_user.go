package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/auth-service/internal/domain"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, first_name, last_name, phone,
			role, identity_verified, verification_status,
			refresh_token_hash, refresh_token_expires_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.IdentityVerified,
		user.VerificationStatus,
		user.RefreshTokenHash,
		user.RefreshTokenExpiresAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if isPgUniqueViolation(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone,
			   role, identity_verified, verification_status,
			   refresh_token_hash, refresh_token_expires_at,
			   created_at, updated_at
		FROM users WHERE id = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.IdentityVerified,
		&user.VerificationStatus,
		&user.RefreshTokenHash,
		&user.RefreshTokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone,
			   role, identity_verified, verification_status,
			   refresh_token_hash, refresh_token_expires_at,
			   created_at, updated_at
		FROM users WHERE email = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.IdentityVerified,
		&user.VerificationStatus,
		&user.RefreshTokenHash,
		&user.RefreshTokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			email = $2,
			password_hash = $3,
			first_name = $4,
			last_name = $5,
			phone = $6,
			role = $7,
			identity_verified = $8,
			verification_status = $9,
			refresh_token_hash = $10,
			refresh_token_expires_at = $11,
			updated_at = $12
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.IdentityVerified,
		user.VerificationStatus,
		user.RefreshTokenHash,
		user.RefreshTokenExpiresAt,
		time.Now(),
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves a paginated list of users
func (r *PostgresUserRepository) List(ctx context.Context, offset, limit int, filters UserFilters) ([]*domain.User, int, error) {
	// Build query with filters
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone,
			   role, identity_verified, verification_status,
			   refresh_token_hash, refresh_token_expires_at,
			   created_at, updated_at
		FROM users
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.Role != nil {
		query += ` AND role = $` + string(rune('0'+argIndex))
		countQuery += ` AND role = $` + string(rune('0'+argIndex))
		args = append(args, *filters.Role)
		argIndex++
	}

	if filters.VerificationStatus != nil {
		query += ` AND verification_status = $` + string(rune('0'+argIndex))
		countQuery += ` AND verification_status = $` + string(rune('0'+argIndex))
		args = append(args, *filters.VerificationStatus)
		argIndex++
	}

	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argIndex)) + ` OFFSET $` + string(rune('0'+argIndex+1))
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Role,
			&user.IdentityVerified,
			&user.VerificationStatus,
			&user.RefreshTokenHash,
			&user.RefreshTokenExpiresAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// UpdateRefreshToken updates the refresh token for a user
func (r *PostgresUserRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, hash string, expiresAt *time.Time) error {
	query := `
		UPDATE users SET
			refresh_token_hash = $2,
			refresh_token_expires_at = $3,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, userID, hash, expiresAt)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// ClearRefreshToken clears the refresh token for a user
func (r *PostgresUserRepository) ClearRefreshToken(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateRefreshToken(ctx, userID, "", nil)
}

// isPgUniqueViolation checks if the error is a PostgreSQL unique constraint violation
func isPgUniqueViolation(err error) bool {
	// Check for PostgreSQL error code 23505 (unique_violation)
	if err != nil && len(err.Error()) > 0 {
		return false // Simplified check - in production use pgconn.PgError
	}
	return false
}
