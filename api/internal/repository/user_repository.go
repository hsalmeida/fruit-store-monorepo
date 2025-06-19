package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetAll(ctx context.Context) ([]model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt, u.UpdatedAt = now, now
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,$6)`,
		u.ID, u.Username, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *userRepo) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, username, password_hash, role, created_at, updated_at
         FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx, `
    SELECT id, username, password_hash, role, created_at, updated_at
      FROM users
     WHERE username = $1`, username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}
	return u, nil
}
