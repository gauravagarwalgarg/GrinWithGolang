package domain

// ErrorResponse is a standard error envelope for all API errors.
type ErrorResponse struct {
	Message string `json:"message"`
}
