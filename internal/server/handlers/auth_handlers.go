package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func LogUserIn(hw *writer.HttpWriter, authSvc user_auth.IAuthService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest user_auth.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			hw.WriteError(w, err)
			return
		}
		if err := r.Body.Close(); err != nil {
			logger.Error("Failed to close request body", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		loginResp, err := authSvc.LogUserIn(r.Context(), &loginRequest)
		if err != nil {
			logger.Error("Failed to log user in", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "User logged in successfully", loginResp)
		hw.WriteResponse(w, http.StatusOK, resp)
		return
	}
}
