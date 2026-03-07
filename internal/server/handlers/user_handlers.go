package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateUserHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[users.CreateUser](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		usr, err := usrSvc.CreateUser(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create user", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("User created successfully", zap.Any("user", usr))
		resp := internal.NewApiResponse(false, "User created successfully", usr)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetUserByIdHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		user, err1 := usrSvc.GetUserById(r.Context(), nil, int32(userId))
		if err1 != nil {
			logger.Error("Failed to get user by ID", zap.Int("user_id", userId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "User retrieved successfully", user)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetUserByEmailHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.PathValue("email")

		user, err1 := usrSvc.GetUserByEmail(r.Context(), nil, email)
		if err1 != nil {
			logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "User retrieved successfully", user)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetUsersHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allUsers, err := usrSvc.GetUsers(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get users", zap.Error(err))
			hw.WriteError(w, err)
		}

		resp := internal.NewApiResponse(false, "Users retrieved successfully", allUsers)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateUserHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
		}

		request, decodeErr := decodeRequestBodyInto[users.UpdateUserRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateUser := &users.UpdateUser{
			UserId:    int32(userId),
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
			Roles:     request.Roles,
		}

		usr, updateErr := usrSvc.UpdateUser(r.Context(), nil, updateUser)
		if updateErr != nil {
			logger.Error("Failed to update user", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "User updated successfully", usr)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateUserPasswordHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		request, decodeErr := decodeRequestBodyInto[users.UpdateUserPasswordRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateUserPass := &users.UpdateUserPassword{
			UserId:   int32(userId),
			Password: request.Password,
		}

		updateErr := usrSvc.UpdateUserPassword(r.Context(), nil, updateUserPass)
		if updateErr != nil {
			logger.Error("Failed to update user password", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "User password updated successfully", nil)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteUserHandler(hw *writer.HttpWriter, usrSvc users.IUserService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		delErr := usrSvc.DeleteUser(r.Context(), nil, int32(userId))
		if delErr != nil {
			logger.Error("Failed to delete user", zap.Int("user_id", userId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "User deleted successfully", nil))
	}
}
