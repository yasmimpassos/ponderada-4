package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"racha-historico/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	query := `
		INSERT INTO users (id, name, email, password_hash, avatar_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash, user.AvatarURL, user.CreatedAt,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, name, email, password_hash, COALESCE(avatar_url, ''), COALESCE(fcm_token, ''), created_at
		FROM users
		WHERE email = ?
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.AvatarURL, &user.FCMToken, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, name, email, COALESCE(avatar_url, ''), COALESCE(fcm_token, ''), created_at
		FROM users
		WHERE id = ?
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email,
		&user.AvatarURL, &user.FCMToken, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateFCMToken(ctx context.Context, userID string, token string) error {
	query := `UPDATE users SET fcm_token = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, token, userID)
	return err
}
