package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"telegram-job/internal/domain"
)

type JobRepository struct {
	db *DB
}

func NewJobRepository(db *DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	query := `
		INSERT INTO jobs (id, company_id, title, level, type, category, salary_from, salary_to, description, apply_link, status, language)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at
	`
	job.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		job.ID,
		job.CompanyID,
		job.Title,
		job.Level,
		job.Type,
		job.Category,
		job.SalaryFrom,
		job.SalaryTo,
		job.Description,
		job.ApplyLink,
		job.Status,
		job.Language,
	).Scan(&job.CreatedAt)
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	query := `
		SELECT id, company_id, title, level, type, category, salary_from, salary_to, description, apply_link, status, language, channel_message_id, published_at, created_at
		FROM jobs
		WHERE id = $1
	`
	var job domain.Job
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.CompanyID,
		&job.Title,
		&job.Level,
		&job.Type,
		&job.Category,
		&job.SalaryFrom,
		&job.SalaryTo,
		&job.Description,
		&job.ApplyLink,
		&job.Status,
		&job.Language,
		&job.ChannelMessageID,
		&job.PublishedAt,
		&job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) GetByStatus(ctx context.Context, status domain.JobStatus) ([]domain.JobWithCompany, error) {
	query := `
		SELECT j.id, j.company_id, j.title, j.level, j.type, j.category, j.salary_from, j.salary_to, j.description, j.apply_link, j.status, j.language, j.published_at, j.created_at, c.name, c.contact, u.telegram_id
		FROM jobs j
		JOIN companies c ON j.company_id = c.id
		JOIN users u ON c.user_id = u.id
		WHERE j.status = $1
		ORDER BY j.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.JobWithCompany
	for rows.Next() {
		var job domain.JobWithCompany
		err := rows.Scan(
			&job.ID,
			&job.CompanyID,
			&job.Title,
			&job.Level,
			&job.Type,
			&job.Category,
			&job.SalaryFrom,
			&job.SalaryTo,
			&job.Description,
			&job.ApplyLink,
			&job.Status,
			&job.Language,
			&job.PublishedAt,
			&job.CreatedAt,
			&job.CompanyName,
			&job.CompanyContact,
			&job.AuthorTelegramID,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) GetWithCompany(ctx context.Context, id uuid.UUID) (*domain.JobWithCompany, error) {
	query := `
		SELECT j.id, j.company_id, j.title, j.level, j.type, j.category, j.salary_from, j.salary_to, j.description, j.apply_link, j.status, j.language, j.published_at, j.created_at, c.name, c.contact, u.telegram_id
		FROM jobs j
		JOIN companies c ON j.company_id = c.id
		JOIN users u ON c.user_id = u.id
		WHERE j.id = $1
	`
	var job domain.JobWithCompany
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.CompanyID,
		&job.Title,
		&job.Level,
		&job.Type,
		&job.Category,
		&job.SalaryFrom,
		&job.SalaryTo,
		&job.Description,
		&job.ApplyLink,
		&job.Status,
		&job.Language,
		&job.PublishedAt,
		&job.CreatedAt,
		&job.CompanyName,
		&job.CompanyContact,
		&job.AuthorTelegramID,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus) error {
	query := `UPDATE jobs SET status = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *JobRepository) SetPublished(ctx context.Context, id uuid.UUID, channelMessageID int) error {
	query := `UPDATE jobs SET status = $1, published_at = $2, channel_message_id = $3 WHERE id = $4`
	_, err := r.db.Pool.Exec(ctx, query, domain.JobStatusPublished, time.Now().UTC(), channelMessageID, id)
	return err
}

func (r *JobRepository) Archive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE jobs SET status = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, domain.JobStatusArchived, id)
	return err
}

func (r *JobRepository) GetExpiredJobs(ctx context.Context, days int) ([]domain.Job, error) {
	query := `
		SELECT id, company_id, title, level, type, category, salary_from, salary_to, description, apply_link, status, language, channel_message_id, published_at, created_at
		FROM jobs
		WHERE status = 'published' AND published_at < NOW() - INTERVAL '1 day' * $1
	`
	rows, err := r.db.Pool.Query(ctx, query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(
			&job.ID,
			&job.CompanyID,
			&job.Title,
			&job.Level,
			&job.Type,
			&job.Category,
			&job.SalaryFrom,
			&job.SalaryTo,
			&job.Description,
			&job.ApplyLink,
			&job.Status,
			&job.Language,
			&job.ChannelMessageID,
			&job.PublishedAt,
			&job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) GetByUserTelegramID(ctx context.Context, telegramID int64) ([]domain.Job, error) {
	query := `
		SELECT j.id, j.company_id, j.title, j.level, j.type, j.category, j.salary_from, j.salary_to, j.description, j.apply_link, j.status, j.language, j.channel_message_id, j.published_at, j.created_at
		FROM jobs j
		JOIN companies c ON j.company_id = c.id
		JOIN users u ON c.user_id = u.id
		WHERE u.telegram_id = $1
		ORDER BY j.created_at DESC
		LIMIT 20
	`
	rows, err := r.db.Pool.Query(ctx, query, telegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(
			&job.ID,
			&job.CompanyID,
			&job.Title,
			&job.Level,
			&job.Type,
			&job.Category,
			&job.SalaryFrom,
			&job.SalaryTo,
			&job.Description,
			&job.ApplyLink,
			&job.Status,
			&job.Language,
			&job.ChannelMessageID,
			&job.PublishedAt,
			&job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) GetStats(ctx context.Context) (*domain.Stats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'published') as published,
			COUNT(*) FILTER (WHERE status = 'rejected') as rejected,
			COUNT(*) FILTER (WHERE status = 'archived') as archived
		FROM jobs
	`
	var stats domain.Stats
	err := r.db.Pool.QueryRow(ctx, query).Scan(
		&stats.Total,
		&stats.Pending,
		&stats.Published,
		&stats.Rejected,
		&stats.Archived,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
