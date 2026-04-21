package apierror

import "net/http"

type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

func New(statusCode int, code, message string) *APIError {
	return &APIError{StatusCode: statusCode, Code: code, Message: message}
}

func NewWithDetails(statusCode int, code, message string, details any) *APIError {
	return &APIError{StatusCode: statusCode, Code: code, Message: message, Details: details}
}

var (
	ErrUnauthorized     = New(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	ErrForbidden        = New(http.StatusForbidden, "FORBIDDEN", "Forbidden")
	ErrNotFound         = New(http.StatusNotFound, "NOT_FOUND", "Resource not found")
	ErrBadRequest       = New(http.StatusBadRequest, "BAD_REQUEST", "Bad request")
	ErrConflict         = New(http.StatusConflict, "CONFLICT", "Resource already exists")
	ErrInternal         = New(http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	ErrInvalidToken     = New(http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
)
