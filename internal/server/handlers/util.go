package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func decodeRequestBodyInto[T any](r *http.Request, logger *zap.Logger) (T, error) {
	var request T

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Failed to decode request body", zap.Error(err))
		return request, fmt.Errorf("failed to decode request body")
	}
	if err := r.Body.Close(); err != nil {
		logger.Error("Failed to close request body", zap.Error(err))
		return request, fmt.Errorf("failed to close request body")
	}

	return request, nil
}
