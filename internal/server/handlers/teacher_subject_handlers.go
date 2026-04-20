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

func LinkSubjectToTeacherHandler(hw *writer.HttpWriter, teacherSubjectSvc teachers.ITeacherSubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := strconv.Atoi(r.PathValue("teacher_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid teacher ID"))
			return
		}

		subjectId, err := strconv.Atoi(r.PathValue("subject_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid subject ID"))
			return
		}

		linkErr := teacherSubjectSvc.LinkSubjectToTeacher(r.Context(), nil, int32(teacherId), int32(subjectId))
		if linkErr != nil {
			logger.Error("Failed to link subject to teacher", zap.Int("teacher_id", teacherId), zap.Int("subject_id", subjectId), zap.Error(linkErr))
			hw.WriteError(w, linkErr)
			return
		}

		logger.Info("Subject linked to teacher successfully", zap.Int("teacher_id", teacherId), zap.Int("subject_id", subjectId))
		resp := internal.NewApiResponse(false, "Subject linked to teacher successfully", nil)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func UnlinkSubjectFromTeacherHandler(hw *writer.HttpWriter, teacherSubjectSvc teachers.ITeacherSubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := strconv.Atoi(r.PathValue("teacher_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid teacher ID"))
			return
		}

		subjectId, err := strconv.Atoi(r.PathValue("subject_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid subject ID"))
			return
		}

		unlinkErr := teacherSubjectSvc.UnlinkSubjectFromTeacher(r.Context(), nil, int32(teacherId), int32(subjectId))
		if unlinkErr != nil {
			logger.Error("Failed to unlink subject from teacher", zap.Int("teacher_id", teacherId), zap.Int("subject_id", subjectId), zap.Error(unlinkErr))
			hw.WriteError(w, unlinkErr)
			return
		}

		logger.Info("Subject unlinked from teacher successfully", zap.Int("teacher_id", teacherId), zap.Int("subject_id", subjectId))
		resp := internal.NewApiResponse(false, "Subject unlinked from teacher successfully", nil)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetSubjectsForTeacherHandler(hw *writer.HttpWriter, teacherSubjectSvc teachers.ITeacherSubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherId, err := strconv.Atoi(r.PathValue("teacher_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid teacher ID"))
			return
		}

		teacherSubjects, getErr := teacherSubjectSvc.GetSubjectsForTeacher(r.Context(), nil, int32(teacherId))
		if getErr != nil {
			logger.Error("Failed to get subjects for teacher", zap.Int("teacher_id", teacherId), zap.Error(getErr))
			hw.WriteError(w, getErr)
			return
		}

		resp := internal.NewApiResponse(false, "Subjects retrieved successfully", teacherSubjects)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

func GetTeachersForSubjectHandler(hw *writer.HttpWriter, teacherSubjectSvc teachers.ITeacherSubjectService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subjectId, err := strconv.Atoi(r.PathValue("subject_id"))
		if err != nil {
			hw.WriteError(w, exceptions.NewRequestValidationError("Invalid subject ID"))
			return
		}

		subjectTeachers, getErr := teacherSubjectSvc.GetTeachersForSubject(r.Context(), nil, int32(subjectId))
		if getErr != nil {
			logger.Error("Failed to get teachers for subject", zap.Int("subject_id", subjectId), zap.Error(getErr))
			hw.WriteError(w, getErr)
			return
		}

		resp := internal.NewApiResponse(false, "Teachers retrieved successfully", subjectTeachers)
		hw.WriteResponse(w, http.StatusOK, resp)
	}
}

