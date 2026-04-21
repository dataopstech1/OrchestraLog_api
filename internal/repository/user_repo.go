package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1 AND is_active = true`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1 AND is_active = true`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, department, avatar_url, is_active)
		VALUES (:id, :email, :password_hash, :first_name, :last_name, :role, :department, :avatar_url, :is_active)
	`, user)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE users SET
			first_name = :first_name,
			last_name = :last_name,
			department = :department,
			avatar_url = :avatar_url,
			updated_at = NOW()
		WHERE id = :id
	`, user)
	return err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET last_login_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]model.User, int, error) {
	var users []model.User
	var total int

	err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users WHERE is_active = true`)
	if err != nil {
		return nil, 0, err
	}

	err = r.db.SelectContext(ctx, &users, `
		SELECT * FROM users WHERE is_active = true
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	return users, total, err
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES (:id, :user_id, :token_hash, :expires_at)
	`, token)
	return err
}

func (r *UserRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.GetContext(ctx, &token, `
		SELECT * FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW()
	`, tokenHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &token, err
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}

func (r *UserRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
	return err
}

func (r *UserRepository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}
