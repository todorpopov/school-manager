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

func LinkParentToStudentHandler(hw *writer.HttpWriter, studentParentSvc students.IStudentParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentId, err := strconv.Atoi(r.PathValue("student_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid student ID"))
			return
		}

		parentId, err := strconv.Atoi(r.PathValue("parent_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid parent ID"))
			return
		}

		linkErr := studentParentSvc.LinkParentToStudent(r.Context(), nil, int32(studentId), int32(parentId))
		if linkErr != nil {
			logger.Error("Failed to link parent to student", zap.Int("student_id", studentId), zap.Int("parent_id", parentId), zap.Error(linkErr))
			hw.WriteError(w, linkErr)
			return
		}

		logger.Info("Parent linked to student successfully", zap.Int("student_id", studentId), zap.Int("parent_id", parentId))
		resp := internal.NewApiResponse(false, "Parent linked to student successfully", nil)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UnlinkParentFromStudentHandler(hw *writer.HttpWriter, studentParentSvc students.IStudentParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentId, err := strconv.Atoi(r.PathValue("student_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid student ID"))
			return
		}

		parentId, err := strconv.Atoi(r.PathValue("parent_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid parent ID"))
			return
		}

		unlinkErr := studentParentSvc.UnlinkParentFromStudent(r.Context(), nil, int32(studentId), int32(parentId))
		if unlinkErr != nil {
			logger.Error("Failed to unlink parent from student", zap.Int("student_id", studentId), zap.Int("parent_id", parentId), zap.Error(unlinkErr))
			hw.WriteError(w, unlinkErr)
			return
		}

		logger.Info("Parent unlinked from student successfully", zap.Int("student_id", studentId), zap.Int("parent_id", parentId))
		resp := internal.NewApiResponse(false, "Parent unlinked from student successfully", nil)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetParentsForStudentHandler(hw *writer.HttpWriter, studentParentSvc students.IStudentParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentId, err := strconv.Atoi(r.PathValue("student_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid student ID"))
			return
		}

		studentParents, getErr := studentParentSvc.GetParentsForStudent(r.Context(), nil, int32(studentId))
		if getErr != nil {
			logger.Error("Failed to get parents for student", zap.Int("student_id", studentId), zap.Error(getErr))
			hw.WriteError(w, getErr)
			return
		}

		resp := internal.NewApiResponse(false, "Parents retrieved successfully", studentParents)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetStudentsForParentHandler(hw *writer.HttpWriter, studentParentSvc students.IStudentParentService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parentId, err := strconv.Atoi(r.PathValue("parent_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid parent ID"))
			return
		}

		parentStudents, getErr := studentParentSvc.GetStudentsForParent(r.Context(), nil, int32(parentId))
		if getErr != nil {
			logger.Error("Failed to get students for parent", zap.Int("parent_id", parentId), zap.Error(getErr))
			hw.WriteError(w, getErr)
			return
		}

		resp := internal.NewApiResponse(false, "Students retrieved successfully", parentStudents)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

