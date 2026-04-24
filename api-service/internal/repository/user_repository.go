package repository

import (
	"context"
	"fmt"

	"github.com/Geze296/orderhub/api-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {

	q := `INSERT INTO users (name, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, created_at
		`
	err := r.db.QueryRow(ctx, q, user.Name, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}


func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT id, name, email, created_at
			FROM users
			WHERE email = $1
		`
	var user domain.User
	err := r.db.QueryRow(ctx, q, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Get by id error:%w", err)
	}

	return &user, nil
}


func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	q := `SELECT id, name, email, created_at
			FROM users
			WHERE id = $1
		`
	var user domain.User
	err := r.db.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Get by Id error: %w", err)
	}

	return &user, nil
}
