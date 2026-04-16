package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateClassHandler(hw *writer.HttpWriter, classSvc classes.IClassService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[classes.CreateClass](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		class, err := classSvc.CreateClass(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create class", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Class created successfully", zap.Any("class", class))
		resp := internal.NewApiResponse(false, "Class created successfully", class)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetClassByIdHandler(hw *writer.HttpWriter, classSvc classes.IClassService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		classId, err := strconv.Atoi(r.PathValue("class_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid class ID"))
			return
		}

		class, err1 := classSvc.GetClassById(r.Context(), nil, int32(classId))
		if err1 != nil {
			logger.Error("Failed to get class by ID", zap.Int("class_id", classId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Class retrieved successfully", class)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetClassesHandler(hw *writer.HttpWriter, classSvc classes.IClassService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allClasses, err := classSvc.GetClasses(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get classes", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Classes retrieved successfully", allClasses)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteClassHandler(hw *writer.HttpWriter, classSvc classes.IClassService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		classId, err := strconv.Atoi(r.PathValue("class_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid class ID"))
			return
		}

		delErr := classSvc.DeleteClass(r.Context(), nil, int32(classId))
		if delErr != nil {
			logger.Error("Failed to delete class", zap.Int("class_id", classId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Class deleted successfully", nil))
	}
}
