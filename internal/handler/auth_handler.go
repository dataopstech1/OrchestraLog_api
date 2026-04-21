package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/middleware"
	"github.com/orchestralog/api/internal/service"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}

	res, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, res)
}

// POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req request.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	_ = h.authService.Logout(r.Context(), req.RefreshToken)
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req request.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	res, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, res)
}

// GET /auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	res, err := h.authService.Me(r.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, res)
}
