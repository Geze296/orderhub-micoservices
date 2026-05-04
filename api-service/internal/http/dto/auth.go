package dto

import "github.com/Geze296/orderhub/api-service/internal/domain"

type RegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}