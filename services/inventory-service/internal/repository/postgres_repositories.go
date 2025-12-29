package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/inventory-service/internal/domain"
)

// PostgresItemRepository implements ItemRepository using PostgreSQL
type PostgresItemRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresItemRepository creates a new PostgreSQL item repository
func NewPostgresItemRepository(pool *pgxpool.Pool) *PostgresItemRepository {
	return &PostgresItemRepository{pool: pool}
}

// Create creates a new rental item
func (r *PostgresItemRepository) Create(ctx context.Context, item *domain.RentalItem) error {
	specsJSON, _ := json.Marshal(item.Specifications)

	query := `
		INSERT INTO rental_items (
			id, owner_id, title, description, category, subcategory,
			daily_rate, weekly_rate, monthly_rate, security_deposit,
			address, city, latitude, longitude,
			specifications, images, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err := r.pool.Exec(ctx, query,
		item.ID, item.OwnerID, item.Title, item.Description, item.Category, item.Subcategory,
		item.DailyRate, item.WeeklyRate, item.MonthlyRate, item.SecurityDeposit,
		item.Address, item.City, item.Latitude, item.Longitude,
		specsJSON, item.Images, item.IsActive, item.CreatedAt, item.UpdatedAt,
	)

	return err
}

// GetByID retrieves an item by ID
func (r *PostgresItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RentalItem, error) {
	query := `
		SELECT id, owner_id, title, description, category, subcategory,
			   daily_rate, weekly_rate, monthly_rate, security_deposit,
			   address, city, latitude, longitude,
			   specifications, images, is_active, created_at, updated_at
		FROM rental_items WHERE id = $1
	`

	item := &domain.RentalItem{}
	var specsJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.OwnerID, &item.Title, &item.Description, &item.Category, &item.Subcategory,
		&item.DailyRate, &item.WeeklyRate, &item.MonthlyRate, &item.SecurityDeposit,
		&item.Address, &item.City, &item.Latitude, &item.Longitude,
		&specsJSON, &item.Images, &item.IsActive, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrItemNotFound
		}
		return nil, err
	}

	if len(specsJSON) > 0 {
		json.Unmarshal(specsJSON, &item.Specifications)
	}

	return item, nil
}

// GetByOwner retrieves items by owner
func (r *PostgresItemRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]*domain.RentalItem, int, error) {
	// Count total
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM rental_items WHERE owner_id = $1", ownerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get items
	query := `
		SELECT id, owner_id, title, description, category, subcategory,
			   daily_rate, weekly_rate, monthly_rate, security_deposit,
			   address, city, latitude, longitude,
			   specifications, images, is_active, created_at, updated_at
		FROM rental_items 
		WHERE owner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanItems(rows, total)
}

// List retrieves items with filters
func (r *PostgresItemRepository) List(ctx context.Context, offset, limit int, filters ItemFilters) ([]*domain.RentalItem, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if filters.Category != nil {
		whereClause += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, *filters.Category)
		argIdx++
	}

	if filters.City != nil {
		whereClause += fmt.Sprintf(" AND city = $%d", argIdx)
		args = append(args, *filters.City)
		argIdx++
	}

	if filters.IsActive != nil {
		whereClause += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filters.IsActive)
		argIdx++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM rental_items " + whereClause
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get items
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT id, owner_id, title, description, category, subcategory,
			   daily_rate, weekly_rate, monthly_rate, security_deposit,
			   address, city, latitude, longitude,
			   specifications, images, is_active, created_at, updated_at
		FROM rental_items 
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanItems(rows, total)
}

// Search searches items by query
func (r *PostgresItemRepository) Search(ctx context.Context, query string, filters ItemFilters, offset, limit int) ([]*domain.RentalItem, int, error) {
	// Simple search implementation
	return r.List(ctx, offset, limit, filters)
}

// Update updates an item
func (r *PostgresItemRepository) Update(ctx context.Context, item *domain.RentalItem) error {
	specsJSON, _ := json.Marshal(item.Specifications)

	query := `
		UPDATE rental_items SET
			title = $2, description = $3, category = $4, subcategory = $5,
			daily_rate = $6, weekly_rate = $7, monthly_rate = $8, security_deposit = $9,
			address = $10, city = $11, latitude = $12, longitude = $13,
			specifications = $14, images = $15, is_active = $16, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		item.ID, item.Title, item.Description, item.Category, item.Subcategory,
		item.DailyRate, item.WeeklyRate, item.MonthlyRate, item.SecurityDeposit,
		item.Address, item.City, item.Latitude, item.Longitude,
		specsJSON, item.Images, item.IsActive,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrItemNotFound
	}

	return nil
}

// Delete deletes an item
func (r *PostgresItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM rental_items WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrItemNotFound
	}

	return nil
}

func (r *PostgresItemRepository) scanItems(rows pgx.Rows, total int) ([]*domain.RentalItem, int, error) {
	items := []*domain.RentalItem{}

	for rows.Next() {
		item := &domain.RentalItem{}
		var specsJSON []byte

		err := rows.Scan(
			&item.ID, &item.OwnerID, &item.Title, &item.Description, &item.Category, &item.Subcategory,
			&item.DailyRate, &item.WeeklyRate, &item.MonthlyRate, &item.SecurityDeposit,
			&item.Address, &item.City, &item.Latitude, &item.Longitude,
			&specsJSON, &item.Images, &item.IsActive, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if len(specsJSON) > 0 {
			json.Unmarshal(specsJSON, &item.Specifications)
		}

		items = append(items, item)
	}

	return items, total, nil
}

// PostgresAvailabilityRepository implements AvailabilityRepository
type PostgresAvailabilityRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAvailabilityRepository(pool *pgxpool.Pool) *PostgresAvailabilityRepository {
	return &PostgresAvailabilityRepository{pool: pool}
}

func (r *PostgresAvailabilityRepository) Create(ctx context.Context, slot *domain.AvailabilitySlot) error {
	query := `
		INSERT INTO availability_slots (id, rental_item_id, start_date, end_date, status, booking_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query, slot.ID, slot.RentalItemID, slot.StartDate, slot.EndDate, slot.Status, slot.BookingID, slot.CreatedAt)
	return err
}

func (r *PostgresAvailabilityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AvailabilitySlot, error) {
	query := `SELECT id, rental_item_id, start_date, end_date, status, booking_id, created_at FROM availability_slots WHERE id = $1`

	slot := &domain.AvailabilitySlot{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&slot.ID, &slot.RentalItemID, &slot.StartDate, &slot.EndDate, &slot.Status, &slot.BookingID, &slot.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSlotNotFound
		}
		return nil, err
	}
	return slot, nil
}

func (r *PostgresAvailabilityRepository) GetByItem(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time) ([]*domain.AvailabilitySlot, error) {
	query := `
		SELECT id, rental_item_id, start_date, end_date, status, booking_id, created_at
		FROM availability_slots
		WHERE rental_item_id = $1 AND start_date >= $2 AND end_date <= $3
		ORDER BY start_date
	`

	rows, err := r.pool.Query(ctx, query, itemID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := []*domain.AvailabilitySlot{}
	for rows.Next() {
		slot := &domain.AvailabilitySlot{}
		err := rows.Scan(&slot.ID, &slot.RentalItemID, &slot.StartDate, &slot.EndDate, &slot.Status, &slot.BookingID, &slot.CreatedAt)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	return slots, nil
}

func (r *PostgresAvailabilityRepository) Update(ctx context.Context, slot *domain.AvailabilitySlot) error {
	query := `UPDATE availability_slots SET status = $2, booking_id = $3 WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, slot.ID, slot.Status, slot.BookingID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrSlotNotFound
	}
	return nil
}

func (r *PostgresAvailabilityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM availability_slots WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrSlotNotFound
	}
	return nil
}

func (r *PostgresAvailabilityRepository) CheckConflict(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time, excludeSlotID *uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM availability_slots
		WHERE rental_item_id = $1
		AND status != 'available'
		AND (start_date, end_date) OVERLAPS ($2, $3)
	`
	args := []interface{}{itemID, startDate, endDate}

	if excludeSlotID != nil {
		query += " AND id != $4"
		args = append(args, *excludeSlotID)
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count > 0, err
}

// PostgresMaintenanceRepository implements MaintenanceRepository
type PostgresMaintenanceRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresMaintenanceRepository(pool *pgxpool.Pool) *PostgresMaintenanceRepository {
	return &PostgresMaintenanceRepository{pool: pool}
}

func (r *PostgresMaintenanceRepository) Create(ctx context.Context, log *domain.MaintenanceLog) error {
	query := `
		INSERT INTO maintenance_logs (id, rental_item_id, maintenance_type, description, start_date, end_date, cost, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query, log.ID, log.RentalItemID, log.MaintenanceType, log.Description, log.StartDate, log.EndDate, log.Cost, log.Status, log.CreatedAt)
	return err
}

func (r *PostgresMaintenanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MaintenanceLog, error) {
	query := `SELECT id, rental_item_id, maintenance_type, description, start_date, end_date, cost, status, created_at FROM maintenance_logs WHERE id = $1`

	log := &domain.MaintenanceLog{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&log.ID, &log.RentalItemID, &log.MaintenanceType, &log.Description, &log.StartDate, &log.EndDate, &log.Cost, &log.Status, &log.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMaintenanceNotFound
		}
		return nil, err
	}
	return log, nil
}

func (r *PostgresMaintenanceRepository) GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.MaintenanceLog, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM maintenance_logs WHERE rental_item_id = $1", itemID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, rental_item_id, maintenance_type, description, start_date, end_date, cost, status, created_at
		FROM maintenance_logs
		WHERE rental_item_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, itemID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	logs := []*domain.MaintenanceLog{}
	for rows.Next() {
		log := &domain.MaintenanceLog{}
		err := rows.Scan(&log.ID, &log.RentalItemID, &log.MaintenanceType, &log.Description, &log.StartDate, &log.EndDate, &log.Cost, &log.Status, &log.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *PostgresMaintenanceRepository) Update(ctx context.Context, log *domain.MaintenanceLog) error {
	query := `UPDATE maintenance_logs SET status = $2, end_date = $3 WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, log.ID, log.Status, log.EndDate)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrMaintenanceNotFound
	}
	return nil
}

func (r *PostgresMaintenanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM maintenance_logs WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrMaintenanceNotFound
	}
	return nil
}
