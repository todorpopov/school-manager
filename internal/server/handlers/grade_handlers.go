package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/grades"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateGradeHandler(hw *writer.HttpWriter, gradeSvc grades.IGradeService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[grades.CreateGrade](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		grade, err := gradeSvc.CreateGrade(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create grade", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Grade created successfully", zap.Any("grade", grade))
		resp := internal.NewApiResponse(false, "Grade created successfully", grade)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetGradeByIdHandler(hw *writer.HttpWriter, gradeSvc grades.IGradeService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gradeId, err := strconv.Atoi(r.PathValue("grade_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid grade ID"))
			return
		}

		grade, err1 := gradeSvc.GetGradeById(r.Context(), nil, int32(gradeId))
		if err1 != nil {
			logger.Error("Failed to get grade by ID", zap.Int("grade_id", gradeId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Grade retrieved successfully", grade)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetGradesHandler(hw *writer.HttpWriter, gradeSvc grades.IGradeService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allGrades, err := gradeSvc.GetGrades(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get grades", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Grades retrieved successfully", allGrades)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteGradeHandler(hw *writer.HttpWriter, gradeSvc grades.IGradeService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gradeId, err := strconv.Atoi(r.PathValue("grade_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid grade ID"))
			return
		}

		delErr := gradeSvc.DeleteGrade(r.Context(), nil, int32(gradeId))
		if delErr != nil {
			logger.Error("Failed to delete grade", zap.Int("grade_id", gradeId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Grade deleted successfully", nil))
	}
}

