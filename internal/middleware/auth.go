package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/response"
)

type contextKey string

const (
	ContextKeyUserID     contextKey = "user_id"
	ContextKeyUserEmail  contextKey = "user_email"
	ContextKeyUserRole   contextKey = "user_role"
	ContextKeyDepartment contextKey = "user_department"
)

type AccessTokenClaims struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	jwt.RegisteredClaims
}

func Auth(accessSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.Error(w, apierror.ErrUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &AccessTokenClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, apierror.ErrInvalidToken
				}
				return []byte(accessSecret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, apierror.ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyUserEmail, claims.Email)
			ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)
			ctx = context.WithValue(ctx, ContextKeyDepartment, claims.Department)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) string {
	v, _ := r.Context().Value(ContextKeyUserID).(string)
	return v
}

func GetUserRole(r *http.Request) string {
	v, _ := r.Context().Value(ContextKeyUserRole).(string)
	return v
}
