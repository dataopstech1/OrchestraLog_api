package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/dto/response"
	"github.com/orchestralog/api/internal/model"
	"github.com/orchestralog/api/internal/repository"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	pkgResponse "github.com/orchestralog/api/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) List(ctx context.Context, p pagination.Params) ([]response.UserResponse, *pkgResponse.Meta, error) {
	users, total, err := s.userRepo.List(ctx, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	out := make([]response.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(&u))
	}
	return out, buildMeta(p, total), nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if user == nil {
		return nil, apierror.ErrNotFound
	}
	r := toUserResponse(user)
	return &r, nil
}

func (s *UserService) Create(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apierror.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierror.ErrInternal
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		Department:   req.Department,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apierror.ErrInternal
	}

	r := toUserResponse(user)
	return &r, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if user == nil {
		return nil, apierror.ErrNotFound
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Department != nil {
		user.Department = req.Department
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apierror.ErrInternal
	}

	r := toUserResponse(user)
	return &r, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return apierror.ErrInternal
	}
	if user == nil {
		return apierror.ErrNotFound
	}
	user.IsActive = false
	return s.userRepo.Update(ctx, user)
}
