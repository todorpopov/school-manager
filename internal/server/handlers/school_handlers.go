package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateSchoolHandler(hw *writer.HttpWriter, schoolSvc schools.ISchoolService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[schools.CreateSchool](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		school, err := schoolSvc.CreateSchool(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create school", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("School created successfully", zap.Any("school", school))
		resp := internal.NewApiResponse(false, "School created successfully", school)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetSchoolByIdHandler(hw *writer.HttpWriter, schoolSvc schools.ISchoolService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, err := strconv.Atoi(r.PathValue("school_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid school ID"))
			return
		}

		school, err1 := schoolSvc.GetSchoolById(r.Context(), nil, int32(schoolId))
		if err1 != nil {
			logger.Error("Failed to get school by ID", zap.Int("school_id", schoolId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "School retrieved successfully", school)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetSchoolsHandler(hw *writer.HttpWriter, schoolSvc schools.ISchoolService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allSchools, err := schoolSvc.GetSchools(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get schools", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Schools retrieved successfully", allSchools)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateSchoolHandler(hw *writer.HttpWriter, schoolSvc schools.ISchoolService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, err := strconv.Atoi(r.PathValue("school_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid school ID"))
			return
		}

		request, decodeErr := decodeRequestBodyInto[schools.UpdateSchoolRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateSchool := &schools.UpdateSchool{
			SchoolId:      int32(schoolId),
			SchoolName:    request.SchoolName,
			SchoolAddress: request.SchoolAddress,
		}

		school, updateErr := schoolSvc.UpdateSchool(r.Context(), nil, updateSchool)
		if updateErr != nil {
			logger.Error("Failed to update school", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "School updated successfully", school)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteSchoolHandler(hw *writer.HttpWriter, schoolSvc schools.ISchoolService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schoolId, err := strconv.Atoi(r.PathValue("school_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid school ID"))
			return
		}

		delErr := schoolSvc.DeleteSchool(r.Context(), nil, int32(schoolId))
		if delErr != nil {
			logger.Error("Failed to delete school", zap.Int("school_id", schoolId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "School deleted successfully", nil))
	}
}

