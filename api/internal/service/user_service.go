package service

import (
	"context"
	"errors"

	"github.com/hsalmeida/fruit-store-monorepo/api/internal/model"
	"github.com/hsalmeida/fruit-store-monorepo/api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetAllUsers() ([]model.User, error)
	Authenticate(ctx context.Context, username, password string) (model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(ctx context.Context, u *model.User) error {
	return s.repo.Create(ctx, u)
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAll(context.Background())
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (model.User, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return u, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return u, ErrInvalidCredentials
	}
	return u, nil
}
