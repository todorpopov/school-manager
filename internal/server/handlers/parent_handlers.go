package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/parents"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateParentHandler(hw *writer.HttpWriter, parentSvc parents.IParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[parents.CreateParent](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		parent, err := parentSvc.CreateParent(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create parent", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Parent created successfully", zap.Any("parent", parent))
		resp := internal.NewApiResponse(false, "Parent created successfully", parent)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetParentByIdHandler(hw *writer.HttpWriter, parentSvc parents.IParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentId, err := strconv.Atoi(r.PathValue("parent_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid parent ID"))
			return
		}

		parent, err1 := parentSvc.GetParentById(r.Context(), nil, int32(parentId))
		if err1 != nil {
			logger.Error("Failed to get parent by ID", zap.Int("parent_id", parentId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Parent retrieved successfully", parent)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetParentByUserIdHandler(hw *writer.HttpWriter, parentSvc parents.IParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		parent, err1 := parentSvc.GetParentByUserId(r.Context(), nil, int32(userId))
		if err1 != nil {
			logger.Error("Failed to get parent by user ID", zap.Int("user_id", userId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Parent retrieved successfully", parent)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetParentsHandler(hw *writer.HttpWriter, parentSvc parents.IParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allParents, err := parentSvc.GetParents(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get parents", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Parents retrieved successfully", allParents)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateParentHandler(hw *writer.HttpWriter, parentSvc parents.IParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentId, err := strconv.Atoi(r.PathValue("parent_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid parent ID"))
			return
		}

		request, decodeErr := decodeRequestBodyInto[parents.UpdateParentRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateParent := &parents.UpdateParent{
			ParentId:  int32(parentId),
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
		}

		parent, updateErr := parentSvc.UpdateParent(r.Context(), nil, updateParent)
		if updateErr != nil {
			logger.Error("Failed to update parent", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "Parent updated successfully", parent)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteParentHandler(hw *writer.HttpWriter, parentSvc parents.IParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentId, err := strconv.Atoi(r.PathValue("parent_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid parent ID"))
			return
		}

		delErr := parentSvc.DeleteParent(r.Context(), nil, int32(parentId))
		if delErr != nil {
			logger.Error("Failed to delete parent", zap.Int("parent_id", parentId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Parent deleted successfully", nil))
	}
}

