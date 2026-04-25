package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateTeacherHandler(hw *writer.HttpWriter, teacherSvc teachers.ITeacherService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[teachers.CreateTeacher](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		teacher, err := teacherSvc.CreateTeacher(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create teacher", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Teacher created successfully", zap.Any("teacher", teacher))
		resp := internal.NewApiResponse(false, "Teacher created successfully", teacher)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetTeacherByIdHandler(hw *writer.HttpWriter, teacherSvc teachers.ITeacherService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := strconv.Atoi(r.PathValue("teacher_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid teacher ID"))
			return
		}

		teacher, err1 := teacherSvc.GetTeacherById(r.Context(), nil, int32(teacherId))
		if err1 != nil {
			logger.Error("Failed to get teacher by ID", zap.Int("teacher_id", teacherId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Teacher retrieved successfully", teacher)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetTeacherByUserIdHandler(hw *writer.HttpWriter, teacherSvc teachers.ITeacherService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		teacher, err1 := teacherSvc.GetTeacherByUserId(r.Context(), nil, int32(userId))
		if err1 != nil {
			logger.Error("Failed to get teacher by user ID", zap.Int("user_id", userId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Teacher retrieved successfully", teacher)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetTeachersHandler(hw *writer.HttpWriter, teacherSvc teachers.ITeacherService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allTeachers, err := teacherSvc.GetTeachers(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get teachers", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Teachers retrieved successfully", allTeachers)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateTeacherHandler(hw *writer.HttpWriter, teacherSvc teachers.ITeacherService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := strconv.Atoi(r.PathValue("teacher_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid teacher ID"))
			return
		}

		request, decodeErr := decodeRequestBodyInto[teachers.UpdateTeacherRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateTeacher := &teachers.UpdateTeacher{
			TeacherId: int32(teacherId),
			SchoolId:  request.SchoolId,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
		}

		teacher, updateErr := teacherSvc.UpdateTeacher(r.Context(), nil, updateTeacher)
		if updateErr != nil {
			logger.Error("Failed to update teacher", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "Teacher updated successfully", teacher)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteTeacherHandler(hw *writer.HttpWriter, teacherSvc teachers.ITeacherService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := strconv.Atoi(r.PathValue("teacher_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid teacher ID"))
			return
		}

		delErr := teacherSvc.DeleteTeacher(r.Context(), nil, int32(teacherId))
		if delErr != nil {
			logger.Error("Failed to delete teacher", zap.Int("teacher_id", teacherId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Teacher deleted successfully", nil))
	}
}

