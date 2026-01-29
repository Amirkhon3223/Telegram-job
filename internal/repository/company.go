package repository

import (
	"context"

	"github.com/google/uuid"
	"telegram-job/internal/domain"
)

type CompanyRepository struct {
	db *DB
}

func NewCompanyRepository(db *DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (id, user_id, name, contact, telegram)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`
	company.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		company.ID,
		company.UserID,
		company.Name,
		company.Contact,
		company.Telegram,
	).Scan(&company.CreatedAt)
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	query := `
		SELECT id, user_id, name, contact, telegram, created_at
		FROM companies
		WHERE id = $1
	`
	var company domain.Company
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&company.ID,
		&company.UserID,
		&company.Name,
		&company.Contact,
		&company.Telegram,
		&company.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Company, error) {
	query := `
		SELECT id, user_id, name, contact, telegram, created_at
		FROM companies
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var company domain.Company
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&company.ID,
		&company.UserID,
		&company.Name,
		&company.Contact,
		&company.Telegram,
		&company.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &company, nil
}
