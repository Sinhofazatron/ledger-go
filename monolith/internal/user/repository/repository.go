package repository

import (
	"context"
	"database/sql"
	"errors"

	"perfect_api/internal/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByEmailNoErrorNotFound(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	CreateTx(ctx context.Context, tx *sql.Tx, u *model.User) error
	UpdateAvatar(ctx context.Context, id string, path string) error
	SoftDelete(ctx context.Context, id string) error
	UpdateVerificationStatus(ctx context.Context, id string, verified bool) error
	UpdateVerificationStatusTx(ctx context.Context, tx *sql.Tx, id string, verified bool) error
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
	GetByOAuth(ctx context.Context, provider, oauthID string) (*model.User, error)
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) Create(ctx context.Context, u *model.User) error {
	query := `INSERT INTO users (id, full_name, email, password_hash) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.FullName, u.Email, u.PasswordHash)
	return err
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, full_name, email, role, password_hash, oauth_provider, oauth_id, avatar_url, is_verified, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`
	u := &model.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.PasswordHash, &u.OAuthProvider, &u.OAuthID, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return u, nil
}

func (r *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, full_name, email, role, password_hash, avatar_url, is_verified, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`
	u := &model.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.PasswordHash, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return u, nil
}

func (r *mysqlUserRepository) GetByEmailNoErrorNotFound(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, full_name, email, role, password_hash, avatar_url, is_verified, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`
	u := &model.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.PasswordHash, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (r *mysqlUserRepository) Update(ctx context.Context, u *model.User) error {
	query := `UPDATE users SET full_name = $1 WHERE id = $2`

	// TODO: Не обновляется updated_at
	_, err := r.db.ExecContext(ctx, query, u.FullName, u.ID)
	return err
}

func (r *mysqlUserRepository) CreateTx(ctx context.Context, tx *sql.Tx, u *model.User) error {
	query := `INSERT INTO users (id, full_name, email, password_hash, oauth_provider, oauth_id, is_verified) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.ExecContext(ctx, query, u.ID, u.FullName, u.Email, u.PasswordHash, u.OAuthProvider, u.OAuthID, u.IsVerified)
	return err
}

func (r *mysqlUserRepository) UpdateAvatar(ctx context.Context, id string, path string) error {
	query := `UPDATE users SET avatar_url = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, path, id)
	return err
}

func (r *mysqlUserRepository) SoftDelete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *mysqlUserRepository) UpdateVerificationStatus(ctx context.Context, id string, verified bool) error {
	query := `UPDATE users SET is_verified = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, verified, id)
	return err
}

func (r *mysqlUserRepository) UpdateVerificationStatusTx(ctx context.Context, tx *sql.Tx, id string, verified bool) error {
	query := `UPDATE users SET is_verified = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := tx.ExecContext(ctx, query, verified, id)
	return err
}

func (r *mysqlUserRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, passwordHash, id)
	return err
}

func (r *mysqlUserRepository) GetByOAuth(ctx context.Context, provider, oauthID string) (*model.User, error) {
	query := `SELECT id, full_name, email, role, password_hash, oauth_provider, oauth_id, avatar_url, is_verified, created_at, updated_at, deleted_at FROM users WHERE oauth_provider = $1 AND oauth_id = $2 AND deleted_at IS NULL`
	u := &model.User{}

	err := r.db.QueryRowContext(ctx, query, provider, oauthID).Scan(
		&u.ID, &u.FullName, &u.Email, &u.Role, &u.PasswordHash, &u.OAuthProvider, &u.OAuthID, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}
