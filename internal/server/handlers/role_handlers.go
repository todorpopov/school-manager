package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/roles"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateRoleHandler(hw *writer.HttpWriter, roleSvc roles.IRoleService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[roles.CreateRole](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		role, err := roleSvc.CreateRole(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create role", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Role created successfully", zap.Any("role", role))
		resp := internal.NewApiResponse(false, "Role created successfully", role)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetRoleByIdHandler(hw *writer.HttpWriter, roleSvc roles.IRoleService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleId, err := strconv.Atoi(r.PathValue("role_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid role ID"))
			return
		}

		role, err1 := roleSvc.GetRoleById(r.Context(), nil, int32(roleId))
		if err1 != nil {
			logger.Error("Failed to get role by ID", zap.Int("role_id", roleId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Role retrieved successfully", role)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetRolesHandler(hw *writer.HttpWriter, roleSvc roles.IRoleService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allRoles, err := roleSvc.GetRoles(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get roles", zap.Error(err))
			hw.WriteError(w, err)
		}

		resp := internal.NewApiResponse(false, "Roles retrieved successfully", allRoles)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteRoleHandler(hw *writer.HttpWriter, roleSvc roles.IRoleService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleId, err := strconv.Atoi(r.PathValue("role_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid role ID"))
			return
		}

		delErr := roleSvc.DeleteRole(r.Context(), nil, int32(roleId))
		if delErr != nil {
			logger.Error("Failed to delete role", zap.Int("role_id", roleId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Role deleted successfully", nil))
	}
}
