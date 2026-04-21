package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/terms"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateTermHandler(hw *writer.HttpWriter, termSvc terms.ITermService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[terms.CreateTerm](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		term, err := termSvc.CreateTerm(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create term", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Term created successfully", zap.Any("term", term))
		resp := internal.NewApiResponse(false, "Term created successfully", term)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetTermByIdHandler(hw *writer.HttpWriter, termSvc terms.ITermService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		termId, err := strconv.Atoi(r.PathValue("term_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid term ID"))
			return
		}

		term, err1 := termSvc.GetTermById(r.Context(), nil, int32(termId))
		if err1 != nil {
			logger.Error("Failed to get term by ID", zap.Int("term_id", termId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Term retrieved successfully", term)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetTermsHandler(hw *writer.HttpWriter, termSvc terms.ITermService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allTerms, err := termSvc.GetTerms(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get terms", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Terms retrieved successfully", allTerms)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteTermHandler(hw *writer.HttpWriter, termSvc terms.ITermService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		termId, err := strconv.Atoi(r.PathValue("term_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid term ID"))
			return
		}

		delErr := termSvc.DeleteTerm(r.Context(), nil, int32(termId))
		if delErr != nil {
			logger.Error("Failed to delete term", zap.Int("term_id", termId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Term deleted successfully", nil))
	}
}

