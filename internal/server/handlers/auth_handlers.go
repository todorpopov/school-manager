package handlers

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterUserHandler(hw *writer.HttpWriter, authSvc user_auth.IAuthService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[user_auth.RegisterRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		request.Roles = []string{"USER"}
		authResp, err := authSvc.RegisterUser(r.Context(), &request)
		if err != nil {
			logger.Error("Failed to register user", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "User registered successfully", authResp)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func RegisterAdminHandler(hw *writer.HttpWriter, authSvc user_auth.IAuthService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[user_auth.RegisterAdminRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		authResp, err := authSvc.RegisterAdminUser(r.Context(), &request)
		if err != nil {
			logger.Error("Failed to register admin user", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Admin user registered successfully", authResp)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func LogUserInHandler(hw *writer.HttpWriter, authSvc user_auth.IAuthService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[user_auth.LoginRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		authResp, err := authSvc.LogUserIn(r.Context(), &request)
		if err != nil {
			logger.Error("Failed to log user in", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "User logged in successfully", authResp)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}
