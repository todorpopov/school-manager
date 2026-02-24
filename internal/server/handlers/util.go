package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/todorpopov/school-manager/internal/exceptions"
	"go.uber.org/zap"
)

func decodeRequestBodyInto[T any](r *http.Request, logger *zap.Logger) (T, *exceptions.AppError) {
	var request T

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Failed to decode request body", zap.Error(err))
		return request, exceptions.NewRequestValidationError("Failed to decode request")
	}
	if err := r.Body.Close(); err != nil {
		logger.Error("Failed to close request body", zap.Error(err))
		return request, exceptions.NewRequestValidationError("Failed to close request body")
	}

	return request, nil
}
