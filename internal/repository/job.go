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

func (r *JobRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts (id, post_type, user_id, company_id, title, level, type, category, salary_from, salary_to, description, apply_link, status, language, experience_years, employment, about, resume_link, contact)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING created_at
	`
	post.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		post.ID,
		post.PostType,
		post.UserID,
		post.CompanyID,
		post.Title,
		post.Level,
		post.Type,
		post.Category,
		post.SalaryFrom,
		post.SalaryTo,
		post.Description,
		post.ApplyLink,
		post.Status,
		post.Language,
		post.ExperienceYears,
		post.Employment,
		post.About,
		post.ResumeLink,
		post.Contact,
	).Scan(&post.CreatedAt)
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	query := `
		SELECT id, post_type, user_id, company_id, title, level, type, category, salary_from, salary_to, description, apply_link, status, language, channel_message_id, published_at, created_at, experience_years, employment, about, resume_link, contact
		FROM posts
		WHERE id = $1
	`
	var post domain.Post
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.PostType,
		&post.UserID,
		&post.CompanyID,
		&post.Title,
		&post.Level,
		&post.Type,
		&post.Category,
		&post.SalaryFrom,
		&post.SalaryTo,
		&post.Description,
		&post.ApplyLink,
		&post.Status,
		&post.Language,
		&post.ChannelMessageID,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.ExperienceYears,
		&post.Employment,
		&post.About,
		&post.ResumeLink,
		&post.Contact,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *JobRepository) GetByStatus(ctx context.Context, status domain.JobStatus) ([]domain.PostWithDetails, error) {
	query := `
		SELECT
			p.id, p.post_type, p.user_id, p.company_id, p.title, p.level, p.type, p.category,
			p.salary_from, p.salary_to, p.description, p.apply_link, p.status, p.language,
			p.published_at, p.created_at, p.experience_years, p.employment, p.about, p.resume_link, p.contact,
			COALESCE(c.name, '') as company_name,
			COALESCE(c.contact, '') as company_contact,
			COALESCE(u.telegram_id, u2.telegram_id) as author_telegram_id
		FROM posts p
		LEFT JOIN companies c ON p.company_id = c.id
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN users u2 ON p.user_id = u2.id
		WHERE p.status = $1
		ORDER BY p.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.PostWithDetails
	for rows.Next() {
		var post domain.PostWithDetails
		err := rows.Scan(
			&post.ID,
			&post.PostType,
			&post.UserID,
			&post.CompanyID,
			&post.Title,
			&post.Level,
			&post.Type,
			&post.Category,
			&post.SalaryFrom,
			&post.SalaryTo,
			&post.Description,
			&post.ApplyLink,
			&post.Status,
			&post.Language,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.ExperienceYears,
			&post.Employment,
			&post.About,
			&post.ResumeLink,
			&post.Contact,
			&post.CompanyName,
			&post.CompanyContact,
			&post.AuthorTelegramID,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *JobRepository) GetWithCompany(ctx context.Context, id uuid.UUID) (*domain.PostWithDetails, error) {
	query := `
		SELECT
			p.id, p.post_type, p.user_id, p.company_id, p.title, p.level, p.type, p.category,
			p.salary_from, p.salary_to, p.description, p.apply_link, p.status, p.language,
			p.published_at, p.created_at, p.experience_years, p.employment, p.about, p.resume_link, p.contact,
			COALESCE(c.name, '') as company_name,
			COALESCE(c.contact, '') as company_contact,
			COALESCE(u.telegram_id, u2.telegram_id) as author_telegram_id
		FROM posts p
		LEFT JOIN companies c ON p.company_id = c.id
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN users u2 ON p.user_id = u2.id
		WHERE p.id = $1
	`
	var post domain.PostWithDetails
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.PostType,
		&post.UserID,
		&post.CompanyID,
		&post.Title,
		&post.Level,
		&post.Type,
		&post.Category,
		&post.SalaryFrom,
		&post.SalaryTo,
		&post.Description,
		&post.ApplyLink,
		&post.Status,
		&post.Language,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.ExperienceYears,
		&post.Employment,
		&post.About,
		&post.ResumeLink,
		&post.Contact,
		&post.CompanyName,
		&post.CompanyContact,
		&post.AuthorTelegramID,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus) error {
	query := `UPDATE posts SET status = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *JobRepository) SetPublished(ctx context.Context, id uuid.UUID, channelMessageID int) error {
	query := `UPDATE posts SET status = $1, published_at = $2, channel_message_id = $3 WHERE id = $4`
	_, err := r.db.Pool.Exec(ctx, query, domain.JobStatusPublished, time.Now().UTC(), channelMessageID, id)
	return err
}

func (r *JobRepository) Archive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE posts SET status = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, domain.JobStatusArchived, id)
	return err
}

func (r *JobRepository) GetExpiredJobs(ctx context.Context, days int) ([]domain.Post, error) {
	query := `
		SELECT id, post_type, user_id, company_id, title, level, type, category, salary_from, salary_to, description, apply_link, status, language, channel_message_id, published_at, created_at, experience_years, employment, about, resume_link, contact
		FROM posts
		WHERE status = 'published' AND published_at < NOW() - INTERVAL '1 day' * $1
	`
	rows, err := r.db.Pool.Query(ctx, query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID,
			&post.PostType,
			&post.UserID,
			&post.CompanyID,
			&post.Title,
			&post.Level,
			&post.Type,
			&post.Category,
			&post.SalaryFrom,
			&post.SalaryTo,
			&post.Description,
			&post.ApplyLink,
			&post.Status,
			&post.Language,
			&post.ChannelMessageID,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.ExperienceYears,
			&post.Employment,
			&post.About,
			&post.ResumeLink,
			&post.Contact,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *JobRepository) GetByUserTelegramID(ctx context.Context, telegramID int64) ([]domain.Post, error) {
	query := `
		SELECT p.id, p.post_type, p.user_id, p.company_id, p.title, p.level, p.type, p.category, p.salary_from, p.salary_to, p.description, p.apply_link, p.status, p.language, p.channel_message_id, p.published_at, p.created_at, p.experience_years, p.employment, p.about, p.resume_link, p.contact
		FROM posts p
		LEFT JOIN companies c ON p.company_id = c.id
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN users u2 ON p.user_id = u2.id
		WHERE u.telegram_id = $1 OR u2.telegram_id = $1
		ORDER BY p.created_at DESC
		LIMIT 20
	`
	rows, err := r.db.Pool.Query(ctx, query, telegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID,
			&post.PostType,
			&post.UserID,
			&post.CompanyID,
			&post.Title,
			&post.Level,
			&post.Type,
			&post.Category,
			&post.SalaryFrom,
			&post.SalaryTo,
			&post.Description,
			&post.ApplyLink,
			&post.Status,
			&post.Language,
			&post.ChannelMessageID,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.ExperienceYears,
			&post.Employment,
			&post.About,
			&post.ResumeLink,
			&post.Contact,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *JobRepository) GetStats(ctx context.Context) (*domain.Stats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'published') as published,
			COUNT(*) FILTER (WHERE status = 'rejected') as rejected,
			COUNT(*) FILTER (WHERE status = 'archived') as archived
		FROM posts
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
