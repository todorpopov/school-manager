package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterTeacherSubjectRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, teacherSubjectSvc teachers.ITeacherSubjectService, authSvc user_auth.IAuthService) {
	logger.Info("Registering teacher-subject relationship routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/teacher-subject/teacher/{teacher_id}/subject/{subject_id}",
		middleware.Chain(
			handlers.LinkSubjectToTeacherHandler(writer, teacherSubjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/teacher-subject/teacher/{teacher_id}/subject/{subject_id}",
		middleware.Chain(
			handlers.UnlinkSubjectFromTeacherHandler(writer, teacherSubjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/teacher-subject/teacher/{teacher_id}/subjects",
		middleware.Chain(
			handlers.GetSubjectsForTeacherHandler(writer, teacherSubjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/teacher-subject/subject/{subject_id}/teachers",
		middleware.Chain(
			handlers.GetTeachersForSubjectHandler(writer, teacherSubjectSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}
