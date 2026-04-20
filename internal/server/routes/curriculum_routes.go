package routes

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/domain_model/curricula"
	"github.com/todorpopov/school-manager/internal/server/handlers"
	"github.com/todorpopov/school-manager/internal/server/middleware"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

func RegisterCurriculumRoutes(s *http.ServeMux, writer *writer.HttpWriter, logger *zap.Logger, curriculumSvc curricula.ICurriculumService, authSvc user_auth.IAuthService) {
	logger.Info("Registering curriculum routes")

	logging := middleware.Logging(logger)
	//requireAdmin := middleware.RequireRoles(writer, authSvc, "ADMIN")

	s.Handle("POST /api/curriculum",
		middleware.Chain(
			handlers.CreateCurriculumHandler(writer, curriculumSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/curriculum/{curriculum_id}",
		middleware.Chain(
			handlers.GetCurriculumByIdHandler(writer, curriculumSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("GET /api/curricula",
		middleware.Chain(
			handlers.GetCurriculaHandler(writer, curriculumSvc, logger),
			logging,
			//requireAdmin,
		),
	)

	s.Handle("DELETE /api/curriculum/{curriculum_id}",
		middleware.Chain(
			handlers.DeleteCurriculumHandler(writer, curriculumSvc, logger),
			logging,
			//requireAdmin,
		),
	)
}

