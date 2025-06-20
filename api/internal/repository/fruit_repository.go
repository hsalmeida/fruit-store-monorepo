package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FruitRepository interface {
	GetAll(ctx context.Context) ([]model.Fruit, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Fruit, error)
	Create(ctx context.Context, f *model.Fruit) error
	Update(ctx context.Context, f *model.Fruit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type fruitRepo struct {
	db *pgxpool.Pool
}

func NewFruitRepository(db *pgxpool.Pool) FruitRepository {
	return &fruitRepo{db: db}
}

func (r *fruitRepo) GetAll(ctx context.Context) ([]model.Fruit, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, quantity, price, created_at FROM fruits`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Fruit
	for rows.Next() {
		var f model.Fruit
		if err := rows.Scan(&f.ID, &f.Name, &f.Quantity, &f.Price, &f.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}

func (r *fruitRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Fruit, error) {
	var f model.Fruit
	err := r.db.QueryRow(ctx, `SELECT id, name, quantity, price, created_at FROM fruits WHERE id=$1`, id).
		Scan(&f.ID, &f.Name, &f.Quantity, &f.Price, &f.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return f, errors.New("fruit not found")
	}
	return f, err
}

func (r *fruitRepo) Create(ctx context.Context, f *model.Fruit) error {
	f.ID = uuid.New()
	f.CreatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`INSERT INTO fruits (id, name, quantity, price, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		f.ID, f.Name, f.Quantity, f.Price, f.CreatedAt, f.UpdatedAt,
	)
	return err
}

func (r *fruitRepo) Update(ctx context.Context, f *model.Fruit) error {
	_, err := r.db.Exec(ctx,
		`UPDATE fruits SET name=$1, quantity=$2, price=$3, updated_at=$4 WHERE id=$5`,
		f.Name, f.Quantity, f.Price, time.Now(), f.ID,
	)
	return err
}

func (r *fruitRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM fruits WHERE id=$1`, id)
	return err
}
