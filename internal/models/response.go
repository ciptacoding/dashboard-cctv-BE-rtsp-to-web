package models

// APIResponse adalah struktur standar untuk response API
type APIResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail memberikan detail error yang lebih spesifik
type ErrorDetail struct {
	Code    string `json:"code"`              // Error code untuk handling di FE
	Message string `json:"message"`           // Human readable message
	Details string `json:"details,omitempty"` // Technical details (optional)
}

// Error codes
const (
	// Authentication errors
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       = "TOKEN_INVALID"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeUserInactive       = "USER_INACTIVE"

	// Validation errors
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeMissingFields    = "MISSING_FIELDS"

	// Resource errors
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"

	// Server errors
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// NewErrorResponse membuat error response dengan error code
func NewErrorResponse(code, message string, details ...string) APIResponse {
	errDetail := &ErrorDetail{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		errDetail.Details = details[0]
	}

	return APIResponse{
		Success: false,
		Message: message,
		Error:   errDetail,
	}
}

// PaginationMeta adalah metadata untuk pagination
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse adalah response dengan pagination
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}
