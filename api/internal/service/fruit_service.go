package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/model"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/repository"
)

type FruitService interface {
	ListFruits(ctx context.Context) ([]model.Fruit, error)
	GetFruit(ctx context.Context, id uuid.UUID) (model.Fruit, error)
	CreateFruit(ctx context.Context, f *model.Fruit) error
	UpdateFruit(ctx context.Context, f *model.Fruit) error
	DeleteFruit(ctx context.Context, id uuid.UUID) error
}

type fruitService struct {
	repo repository.FruitRepository
}

func NewFruitService(r repository.FruitRepository) FruitService {
	return &fruitService{repo: r}
}

func (s *fruitService) ListFruits(ctx context.Context) ([]model.Fruit, error) {
	return s.repo.GetAll(ctx)
}

func (s *fruitService) GetFruit(ctx context.Context, id uuid.UUID) (model.Fruit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *fruitService) CreateFruit(ctx context.Context, f *model.Fruit) error {
	return s.repo.Create(ctx, f)
}

func (s *fruitService) UpdateFruit(ctx context.Context, f *model.Fruit) error {
	return s.repo.Update(ctx, f)
}

func (s *fruitService) DeleteFruit(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
