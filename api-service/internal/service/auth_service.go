package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Geze296/orderhub/api-service/internal/auth"
	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/Geze296/orderhub/api-service/internal/repository"
)

type AuthService struct {
	repo     repository.UserRepository
	jwtSecret string
	ttl       time.Duration
}

var (
	ErrInvalidInput     = errors.New("Invalid input or miss required input")
	ErrShortPasswordLen = errors.New("Short Password length")
	ErrUserExists       = errors.New("User already exist")
	ErrPasswordNotMatch = errors.New("Password Do not match")
)

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:     repo,
		jwtSecret: jwtSecret,
		ttl:       time.Hour,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	if len(input.Password) < 4 {
		return nil, ErrShortPasswordLen
	}

	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, ErrUserExists
	}

	passwordHash, err := auth.HashedPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	e := s.repo.Create(ctx, newUser)

	if e != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(s.jwtSecret, int(newUser.ID), s.ttl)
	if err != nil {
		return nil, fmt.Errorf("Generate token:%w", err)
	}
	return &AuthResult{
		Token: token,
		User:  &domain.User{
			ID: newUser.ID,
			Name: newUser.Name,
			Email: newUser.Email,
			CreatedAt: newUser.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	input.Email = strings.TrimSpace(input.Email)

	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	e := auth.CheckHashedPassword(input.Password, user.PasswordHash)
	if e != nil {
		return nil, ErrPasswordNotMatch
	}

	token, err := auth.GenerateToken(s.jwtSecret, int(user.ID), s.ttl)
	if err != nil {
		return nil, fmt.Errorf("Generate token error:%w", err)
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}
