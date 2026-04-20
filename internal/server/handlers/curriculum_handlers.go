package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/curricula"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateCurriculumHandler(hw *writer.HttpWriter, curriculumSvc curricula.ICurriculumService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[curricula.CreateCurriculum](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		curriculum, err := curriculumSvc.CreateCurriculum(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create curriculum", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Curriculum created successfully", zap.Any("curriculum", curriculum))
		resp := internal.NewApiResponse(false, "Curriculum created successfully", curriculum)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetCurriculumByIdHandler(hw *writer.HttpWriter, curriculumSvc curricula.ICurriculumService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		curriculumId, err := strconv.Atoi(r.PathValue("curriculum_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid curriculum ID"))
			return
		}

		curriculum, err1 := curriculumSvc.GetCurriculumById(r.Context(), nil, int32(curriculumId))
		if err1 != nil {
			logger.Error("Failed to get curriculum by ID", zap.Int("curriculum_id", curriculumId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Curriculum retrieved successfully", curriculum)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetCurriculaHandler(hw *writer.HttpWriter, curriculumSvc curricula.ICurriculumService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allCurricula, err := curriculumSvc.GetCurricula(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get curricula", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Curricula retrieved successfully", allCurricula)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteCurriculumHandler(hw *writer.HttpWriter, curriculumSvc curricula.ICurriculumService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		curriculumId, err := strconv.Atoi(r.PathValue("curriculum_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid curriculum ID"))
			return
		}

		delErr := curriculumSvc.DeleteCurriculum(r.Context(), nil, int32(curriculumId))
		if delErr != nil {
			logger.Error("Failed to delete curriculum", zap.Int("curriculum_id", curriculumId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Curriculum deleted successfully", nil))
	}
}

