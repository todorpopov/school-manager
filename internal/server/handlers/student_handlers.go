package handlers

import (
	"net/http"
	"strconv"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/domain_model/students"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func CreateStudentHandler(hw *writer.HttpWriter, studentSvc students.IStudentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[students.CreateStudent](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		student, err := studentSvc.CreateStudent(r.Context(), nil, &request)
		if err != nil {
			logger.Error("Failed to create student", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Student created successfully", zap.Any("student", student))
		resp := internal.NewApiResponse(false, "Student created successfully", student)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

func GetStudentByIdHandler(hw *writer.HttpWriter, studentSvc students.IStudentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentId, err := strconv.Atoi(r.PathValue("student_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid student ID"))
			return
		}

		student, err1 := studentSvc.GetStudentById(r.Context(), nil, int32(studentId))
		if err1 != nil {
			logger.Error("Failed to get student by ID", zap.Int("student_id", studentId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Student retrieved successfully", student)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetStudentByUserIdHandler(hw *writer.HttpWriter, studentSvc students.IStudentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("user_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid user ID"))
			return
		}

		student, err1 := studentSvc.GetStudentByUserId(r.Context(), nil, int32(userId))
		if err1 != nil {
			logger.Error("Failed to get student by user ID", zap.Int("user_id", userId), zap.Error(err1))
			hw.WriteError(w, err1)
			return
		}

		resp := internal.NewApiResponse(false, "Student retrieved successfully", student)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetStudentsHandler(hw *writer.HttpWriter, studentSvc students.IStudentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allStudents, err := studentSvc.GetStudents(r.Context(), nil)
		if err != nil {
			logger.Error("Failed to get students", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		resp := internal.NewApiResponse(false, "Students retrieved successfully", allStudents)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UpdateStudentHandler(hw *writer.HttpWriter, studentSvc students.IStudentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentId, err := strconv.Atoi(r.PathValue("student_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid student ID"))
			return
		}

		request, decodeErr := decodeRequestBodyInto[students.UpdateStudentRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		updateStudent := &students.UpdateStudent{
			StudentId: int32(studentId),
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
			ClassId:   request.ClassId,
		}

		student, updateErr := studentSvc.UpdateStudent(r.Context(), nil, updateStudent)
		if updateErr != nil {
			logger.Error("Failed to update student", zap.Error(updateErr))
			hw.WriteError(w, updateErr)
			return
		}

		resp := internal.NewApiResponse(false, "Student updated successfully", student)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func DeleteStudentHandler(hw *writer.HttpWriter, studentSvc students.IStudentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentId, err := strconv.Atoi(r.PathValue("student_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid student ID"))
			return
		}

		delErr := studentSvc.DeleteStudent(r.Context(), nil, int32(studentId))
		if delErr != nil {
			logger.Error("Failed to delete student", zap.Int("student_id", studentId), zap.Error(delErr))
			hw.WriteError(w, delErr)
			return
		}

		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "Student deleted successfully", nil))
	}
}

