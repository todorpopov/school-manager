package exceptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal"
	"go.uber.org/zap"
)

type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Cause   error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func NewAppError(code string, message string, cause error) *AppError {
	return &AppError{code, message, nil, cause}
}

func NewValidationError(message string, data map[string]string) *AppError {
	return &AppError{"VALIDATION_ERROR", message, data, nil}
}

func NewRequestValidationError(message string) *AppError {
	return &AppError{"REQUEST_VALIDATION_ERROR", message, nil, nil}
}

type ErrorWriter struct {
	logger *zap.Logger
}

func NewErrorWriter() *ErrorWriter {
	return &ErrorWriter{}
}

func (eh *ErrorWriter) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, data ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := internal.ApiResponse{
		Error:   true,
		Message: message,
		Data:    data,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (eh *ErrorWriter) WriteError(w http.ResponseWriter, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		eh.writeAppError(w, appErr)
		return
	}

	eh.sendErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
}

func (eh *ErrorWriter) writeAppError(w http.ResponseWriter, appErr *AppError) {
	switch appErr.Code {
	case "INTERNAL_ERROR":
		eh.sendErrorResponse(w, http.StatusInternalServerError, appErr.Message)
	case "VALIDATION_ERROR":
		eh.sendErrorResponse(w, http.StatusBadRequest, appErr.Message, appErr.Data)
	case "REQUEST_VALIDATION_ERROR":
		eh.sendErrorResponse(w, http.StatusBadRequest, appErr.Message)
	case "INVALID_CREDENTIALS":
		eh.sendErrorResponse(w, http.StatusUnauthorized, appErr.Message)
	case "NOT_FOUND":
		eh.sendErrorResponse(w, http.StatusNotFound, appErr.Message)
	case "INTEGRITY_CONSTRAINT_VIOLATION":
		eh.sendErrorResponse(w, http.StatusBadRequest, appErr.Message)
	case "RESTRICT_VIOLATION":
		eh.sendErrorResponse(w, http.StatusInternalServerError, appErr.Message)
	case "NOT_NULL_VIOLATION":
		eh.sendErrorResponse(w, http.StatusBadRequest, appErr.Message)
	case "FOREIGN_KEY_VIOLATION":
		eh.sendErrorResponse(w, http.StatusBadRequest, appErr.Message)
	case "UNIQUE_VIOLATION":
		eh.sendErrorResponse(w, http.StatusBadRequest, appErr.Message)
	case "DATABASE_ERROR":
		eh.sendErrorResponse(w, http.StatusInternalServerError, appErr.Message)
	default:
		eh.sendErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
}

func PgErrorToAppError(err error) *AppError {
	appErr := AppError{Code: "DATABASE_ERROR", Message: "Database Error", Cause: err}

	if errors.Is(err, pgx.ErrNoRows) {
		appErr.Code = "NOT_FOUND"
		appErr.Message = "Not found"
		return &appErr
	}

	var pgxErr *pgconn.PgError
	if !errors.As(err, &pgxErr) {
		return &appErr
	}

	switch pgxErr.Code {
	// Class 23 - Integrity constraint violation
	case "23000":
		appErr.Code = "INTEGRITY_CONSTRAINT_VIOLATION"
		appErr.Message = "Integrity constraint violation"
		return &appErr
	case "23001":
		appErr.Code = "RESTRICT_VIOLATION"
		appErr.Message = "Foreign key violation"
		return &appErr
	case "23002":
		appErr.Code = "NOT_NULL_VIOLATION"
		appErr.Message = "Not null violation"
		return &appErr
	case "23503":
		appErr.Code = "FOREIGN_KEY_VIOLATION"
		appErr.Message = "Foreign key violation"
		return &appErr
	case "23505":
		appErr.Code = "UNIQUE_VIOLATION"
		appErr.Message = "Unique violation"
		return &appErr
	default:
		return &appErr
	}
}
