/*
Module 8: Production - Error Handling Patterns

Demonstrates:
  - Structured errors with context (wrapping with %w)
  - Error codes for API responses (machine-readable)
  - Stack context without external packages
  - Error handling middleware pattern
  - Logging errors with request context (ID, user, operation)
  - Sentinel errors vs typed errors vs wrapped errors

Key insight: Production errors need TWO audiences:
  1. Users: clean message + error code (no internals leaked)
  2. Developers: full context, stack trace, request correlation

Run: go run main.go
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

// --- Error types with codes ---

type ErrorCode string

const (
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
)

// AppError is a structured application error with context.
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Op      string    `json:"-"` // operation that failed (for logs)
	Err     error     `json:"-"` // underlying error (for logs)
	Caller  string    `json:"-"` // file:line where error originated
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s: %v", e.Code, e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Op, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

// NewAppError creates an error with automatic caller info.
func NewAppError(code ErrorCode, op, message string, err error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Code:    code,
		Message: message,
		Op:      op,
		Err:     err,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}
}

// --- Request context for error logging ---

type RequestContext struct {
	RequestID string
	UserID    string
	Method    string
	Path      string
	StartTime time.Time
}

// logError logs structured error information for debugging.
func logError(reqCtx *RequestContext, err error) {
	appErr, ok := err.(*AppError)
	if !ok {
		log.Printf("[ERROR] request_id=%s user=%s path=%s error=%v",
			reqCtx.RequestID, reqCtx.UserID, reqCtx.Path, err)
		return
	}
	log.Printf("[ERROR] request_id=%s user=%s op=%s code=%s caller=%s error=%v",
		reqCtx.RequestID, reqCtx.UserID, appErr.Op, appErr.Code,
		appErr.Caller, appErr.Err)
}

// --- Error response middleware ---

func errorResponse(w http.ResponseWriter, err error) {
	appErr, ok := err.(*AppError)
	if !ok {
		// Unknown error: don't leak internals
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    string(ErrInternal),
			"message": "an unexpected error occurred",
		})
		return
	}

	status := mapErrorToHTTP(appErr.Code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"code":    string(appErr.Code),
		"message": appErr.Message, // user-safe message only
	})
}

func mapErrorToHTTP(code ErrorCode) int {
	switch code {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrValidation:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// --- Example service demonstrating error propagation ---

func getUserByID(id string) (*struct{ Name string }, error) {
	if id == "" {
		return nil, NewAppError(ErrValidation, "getUserByID", "user ID is required", nil)
	}
	if id != "1" {
		return nil, NewAppError(ErrNotFound, "getUserByID",
			fmt.Sprintf("user %s not found", id), nil)
	}
	return &struct{ Name string }{"Alice"}, nil
}

func main() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		reqCtx := &RequestContext{
			RequestID: fmt.Sprintf("req_%d", time.Now().UnixNano()),
			UserID:    "anonymous",
			Method:    r.Method,
			Path:      r.URL.Path,
			StartTime: time.Now(),
		}

		id := r.URL.Query().Get("id")
		user, err := getUserByID(id)
		if err != nil {
			logError(reqCtx, err)
			errorResponse(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	log.Println("Error handling demo on :8080")
	log.Println("Try: curl localhost:8080/user?id=1  (success)")
	log.Println("Try: curl localhost:8080/user?id=99 (not found)")
	log.Println("Try: curl localhost:8080/user       (validation)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
