package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/auth-service/internal/domain"
)

// PostgresDocumentRepository implements DocumentRepository using PostgreSQL
type PostgresDocumentRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresDocumentRepository creates a new PostgreSQL document repository
func NewPostgresDocumentRepository(pool *pgxpool.Pool) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{pool: pool}
}

// Create creates a new identity document
func (r *PostgresDocumentRepository) Create(ctx context.Context, doc *domain.IdentityDocument) error {
	query := `
		INSERT INTO identity_documents (id, user_id, document_type, document_url, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		doc.ID,
		doc.UserID,
		doc.DocumentType,
		doc.DocumentURL,
		doc.UploadedAt,
	)

	return err
}

// GetByID retrieves a document by ID
func (r *PostgresDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.IdentityDocument, error) {
	query := `
		SELECT id, user_id, document_type, document_url, uploaded_at
		FROM identity_documents WHERE id = $1
	`

	doc := &domain.IdentityDocument{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&doc.ID,
		&doc.UserID,
		&doc.DocumentType,
		&doc.DocumentURL,
		&doc.UploadedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDocumentNotFound
		}
		return nil, err
	}

	return doc, nil
}

// GetByUserID retrieves all documents for a user
func (r *PostgresDocumentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.IdentityDocument, error) {
	query := `
		SELECT id, user_id, document_type, document_url, uploaded_at
		FROM identity_documents WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := []*domain.IdentityDocument{}
	for rows.Next() {
		doc := &domain.IdentityDocument{}
		err := rows.Scan(
			&doc.ID,
			&doc.UserID,
			&doc.DocumentType,
			&doc.DocumentURL,
			&doc.UploadedAt,
		)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// Delete deletes a document
func (r *PostgresDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM identity_documents WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrDocumentNotFound
	}

	return nil
}
