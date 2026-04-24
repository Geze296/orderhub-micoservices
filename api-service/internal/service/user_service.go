package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	UserRepo *pgxpool.Pool
}