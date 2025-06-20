package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hsalmeida/fruit-store-monorepo/user-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Save(ctx context.Context, u model.User) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Save(ctx context.Context, u model.User) error {

	_, err := r.db.Exec(ctx, `
        INSERT INTO users (id, username, role, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)
    `, uuid.New(), u.Username, u.Role, u.CreatedAt, u.UpdatedAt)
	return err
}
