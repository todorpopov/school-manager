package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/directors"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateDirectorHandler(hw *writer.HttpWriter, directorSvc directors.IDirectorService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[directors.CreateDirector](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		director, err := directorSvc.CreateDirector(r.Context(), &request)
		if err != nil {
			logger.Error("Failed to create director", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Director created successfully", zap.Any("director", director))
		resp := internal.NewApiResponse(false, "Director created successfully", director)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetDirectorByIdHandler(hw *writer.HttpWriter, directorSvc directors.IDirectorService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		directorId, err := strconv.Atoi(r.PathValue("director_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid director ID"))
			return
		}

		director, err1 := directorSvc.GetDirectorById(r.Context(), int32(directorId))
		if err1 != nil {
			logger.Error("Failed to get director by ID", zap.Int("director_id", directorId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Director retrieved successfully", director)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetDirectorByUserIdHandler(hw *writer.HttpWriter, directorSvc directors.IDirectorService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		director, err1 := directorSvc.GetDirectorByUserId(r.Context(), int32(userId))
		if err1 != nil {
			logger.Error("Failed to get director by user ID", zap.Int("user_id", userId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Director retrieved successfully", director)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetDirectorsHandler(hw *writer.HttpWriter, directorSvc directors.IDirectorService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allDirectors, err := directorSvc.GetDirectors(r.Context())
		if err != nil {
			logger.Error("Failed to get directors", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Directors retrieved successfully", allDirectors)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateDirectorHandler(hw *writer.HttpWriter, directorSvc directors.IDirectorService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		directorId, err := strconv.Atoi(r.PathValue("director_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid director ID"))
			return
		}

		request, decodeErr := decodeRequestBodyInto[directors.UpdateDirectorRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateDirector := &directors.UpdateDirector{
			DirectorId: int32(directorId),
			FirstName:  request.FirstName,
			LastName:   request.LastName,
			Email:      request.Email,
		}

		director, updateErr := directorSvc.UpdateDirector(r.Context(), updateDirector)
		if updateErr != nil {
			logger.Error("Failed to update director", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "Director updated successfully", director)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteDirectorHandler(hw *writer.HttpWriter, directorSvc directors.IDirectorService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		directorId, err := strconv.Atoi(r.PathValue("director_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid director ID"))
			return
		}

		delErr := directorSvc.DeleteDirector(r.Context(), int32(directorId))
		if delErr != nil {
			logger.Error("Failed to delete director", zap.Int("director_id", directorId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Director deleted successfully", nil))
	}
}

