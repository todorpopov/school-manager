package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateSubjectHandler(hw *writer.HttpWriter, subjectSvc subjects.ISubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[subjects.CreateSubject](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		subject, err := subjectSvc.CreateSubject(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create subject", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Subject created successfully", zap.Any("subject", subject))
		resp := internal.NewApiResponse(false, "Subject created successfully", subject)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetSubjectByIdHandler(hw *writer.HttpWriter, subjectSvc subjects.ISubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subjectId, err := strconv.Atoi(r.PathValue("subject_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid subject ID"))
			return
		}

		subject, err1 := subjectSvc.GetSubjectById(r.Context(), nil, int32(subjectId))
		if err1 != nil {
			logger.Error("Failed to get subject by ID", zap.Int("subject_id", subjectId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Subject retrieved successfully", subject)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetSubjectsHandler(hw *writer.HttpWriter, subjectSvc subjects.ISubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allSubjects, err := subjectSvc.GetSubjects(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get subjects", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Subjects retrieved successfully", allSubjects)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteSubjectHandler(hw *writer.HttpWriter, subjectSvc subjects.ISubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subjectId, err := strconv.Atoi(r.PathValue("subject_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid subject ID"))
			return
		}

		delErr := subjectSvc.DeleteSubject(r.Context(), nil, int32(subjectId))
		if delErr != nil {
			logger.Error("Failed to delete subject", zap.Int("subject_id", subjectId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Subject deleted successfully", nil))
	}
}

