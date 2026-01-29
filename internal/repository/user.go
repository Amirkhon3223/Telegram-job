package repository

import (
	"context"

	"github.com/google/uuid"
	"telegram-job/internal/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, telegram_id, username, role)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	user.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		user.ID,
		user.TelegramID,
		user.Username,
		user.Role,
	).Scan(&user.CreatedAt)
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `
		SELECT id, telegram_id, username, role, interface_language, created_at
		FROM users
		WHERE telegram_id = $1
	`
	var user domain.User
	err := r.db.Pool.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.Role,
		&user.InterfaceLanguage,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SetInterfaceLanguage(ctx context.Context, telegramID int64, lang string) error {
	query := `UPDATE users SET interface_language = $1 WHERE telegram_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, lang, telegramID)
	return err
}

func (r *UserRepository) GetOrCreate(ctx context.Context, telegramID int64, username string) (*domain.User, error) {
	user, err := r.GetByTelegramID(ctx, telegramID)
	if err == nil {
		return user, nil
	}

	newUser := &domain.User{
		TelegramID: telegramID,
		Username:   username,
		Role:       domain.UserRoleRecruiter,
	}
	if err := r.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
