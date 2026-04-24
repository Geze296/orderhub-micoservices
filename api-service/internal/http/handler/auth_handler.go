package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Geze296/orderhub/api-service/internal/http/helper"
	"github.com/Geze296/orderhub/api-service/internal/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler{
	return &AuthHandler{
		AuthService: authService,
	}
}

type RegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Error while Decoding the request")
	}

	result, err := h.AuthService.Register(r.Context(), service.RegisterInput{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			helper.WriteError(w, http.StatusBadRequest, service.ErrInvalidInput.Error())
		case errors.Is(err, service.ErrShortPasswordLen):
			helper.WriteError(w, http.StatusBadRequest, service.ErrShortPasswordLen.Error())
		case errors.Is(err, service.ErrUserExists):
			helper.WriteError(w, http.StatusBadRequest, service.ErrUserExists.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	helper.WriteJson(w, result, http.StatusCreated, "User Created Successfully")

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Unable to process the request")
	}

	result, err := h.AuthService.Login(r.Context(), service.LoginInput{
		Email: req.Email,
		Password: req.Password,
	})


	if err != nil {
		switch{
		case errors.Is(err, service.ErrInvalidInput):
			helper.WriteError(w, http.StatusBadRequest, service.ErrInvalidInput.Error())
		case errors.Is(err, service.ErrPasswordNotMatch):
			helper.WriteError(w, http.StatusBadRequest, service.ErrPasswordNotMatch.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			fmt.Println(err)
		}
		return
	}

	helper.WriteJson(w, result, http.StatusOK, "Login Successfully!!")
}