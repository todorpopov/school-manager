package exceptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"go.uber.org/zap"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
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
	return &AppError{code, message, cause}
}

type ErrorWriter struct {
	logger *zap.Logger
}

func NewErrorWriter() *ErrorWriter {
	return &ErrorWriter{}
}

func (eh *ErrorWriter) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := internal.ApiResponse{
		Error:   true,
		Message: message,
		Data:    nil,
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
	case "TEST_ERROR":
		eh.sendErrorResponse(w, http.StatusInternalServerError, appErr.Message)
	default:
		eh.sendErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
