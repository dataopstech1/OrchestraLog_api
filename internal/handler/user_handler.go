package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/service"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	"github.com/orchestralog/api/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService, validate: validator.New()}
}

func (h *UserHandler) apiErr(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		response.Error(w, apiErr)
		return
	}
	response.Error(w, apierror.ErrInternal)
}

// GET /users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	p := pagination.Parse(r)
	users, meta, err := h.userService.List(r.Context(), p)
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, users, meta)
}

// GET /users/:id
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(parseURLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	user, svcErr := h.userService.GetByID(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	user, svcErr := h.userService.Create(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, user)
}

// PUT /users/:id
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(parseURLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	user, svcErr := h.userService.Update(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// DELETE /users/:id
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(parseURLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.userService.Delete(r.Context(), id); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "user deactivated"})
}
