package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/absences"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func BulkCreateAbsencesHandler(hw *writer.HttpWriter, absenceSvc absences.IAbsenceService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[absences.BulkCreateAbsences](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		result, err := absenceSvc.BulkCreateAbsences(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to bulk create absences", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Absences bulk created successfully", zap.Int("count", len(result)))
		resp := internal.NewApiResponse(false, "Absences created successfully", result)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func CreateAbsenceHandler(hw *writer.HttpWriter, absenceSvc absences.IAbsenceService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[absences.CreateAbsence](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		absence, err := absenceSvc.CreateAbsence(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create absence", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Absence created successfully", zap.Any("absence", absence))
		resp := internal.NewApiResponse(false, "Absence created successfully", absence)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetAbsenceByIdHandler(hw *writer.HttpWriter, absenceSvc absences.IAbsenceService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		absenceId, err := strconv.Atoi(r.PathValue("absence_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid absence ID"))
			return
		}

		absence, err1 := absenceSvc.GetAbsenceById(r.Context(), nil, int32(absenceId))
		if err1 != nil {
			logger.Error("Failed to get absence by ID", zap.Int("absence_id", absenceId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Absence retrieved successfully", absence)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetAbsencesHandler(hw *writer.HttpWriter, absenceSvc absences.IAbsenceService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAbsences, err := absenceSvc.GetAbsences(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get absences", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Absences retrieved successfully", allAbsences)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteAbsenceHandler(hw *writer.HttpWriter, absenceSvc absences.IAbsenceService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		absenceId, err := strconv.Atoi(r.PathValue("absence_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid absence ID"))
			return
		}

		delErr := absenceSvc.DeleteAbsence(r.Context(), nil, int32(absenceId))
		if delErr != nil {
			logger.Error("Failed to delete absence", zap.Int("absence_id", absenceId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Absence deleted successfully", nil))
	}
}

