package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateUserHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createUser users.CreateUser

		if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			hw.WriteError(w, err)
			return
		}
		if err := r.Body.Close(); err != nil {
			logger.Error("Failed to close request body", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		usr, err := usrSvc.CreateUser(r.Context(), nil, &createUser)
		if err != nil {
			logger.Error("Failed to create user", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("User created successfully", zap.Any("user", usr))
		resp := internal.NewApiResponse(false, "User created successfully", usr)
		hw.WriteResponse(w, http.StatusCreated, resp)
		return
	}
}
