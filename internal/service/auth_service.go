package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/orchestralog/api/internal/config"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/dto/response"
	"github.com/orchestralog/api/internal/model"
	"github.com/orchestralog/api/internal/repository"
	"github.com/orchestralog/api/pkg/apierror"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if user == nil {
		return nil, apierror.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apierror.ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, apierror.ErrInternal
	}

	refreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, apierror.ErrInternal
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWT.AccessExpMinutes * 60,
		User:         toUserResponse(user),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.userRepo.DeleteRefreshToken(ctx, hash)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*response.TokenResponse, error) {
	hash := hashToken(refreshTokenStr)
	stored, err := s.userRepo.FindRefreshToken(ctx, hash)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if stored == nil {
		return nil, apierror.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, stored.UserID)
	if err != nil || user == nil {
		return nil, apierror.ErrInvalidToken
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, apierror.ErrInternal
	}

	return &response.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.cfg.JWT.AccessExpMinutes * 60,
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*response.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, apierror.ErrBadRequest
	}

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

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"role":       user.Role,
		"department": user.Department,
		"exp":        time.Now().Add(time.Duration(s.cfg.JWT.AccessExpMinutes) * time.Minute).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.AccessSecret))
}

func (s *AuthService) generateRefreshToken(ctx context.Context, user *model.User) (string, error) {
	raw := fmt.Sprintf("%s-%s-%d", user.ID, uuid.New().String(), time.Now().UnixNano())
	hash := hashToken(raw)

	rt := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().AddDate(0, 0, s.cfg.JWT.RefreshExpDays),
	}

	if err := s.userRepo.SaveRefreshToken(ctx, rt); err != nil {
		return "", err
	}
	return raw, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func toUserResponse(u *model.User) response.UserResponse {
	return response.UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Role:        u.Role,
		Department:  u.Department,
		AvatarURL:   u.AvatarURL,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
}
