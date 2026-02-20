package server

import (
	"net/http"

	"github.com/todorpopov/school-manager/configs"
	"github.com/todorpopov/school-manager/internal/server/routes"
	"go.uber.org/zap"
)

type Server interface {
	Start()
	RegisterRoutes()
}

type HttpServer struct {
	env    *configs.Config
	mux    *http.ServeMux
	logger *zap.Logger
}

func NewHttpServer(env *configs.Config, logger *zap.Logger) *HttpServer {
	return &HttpServer{
		env,
		http.NewServeMux(),
		logger,
	}
}

func (s *HttpServer) Start() error {
	s.logger.Info("Starting HTTP server on port: " + s.env.ApiPort)
	srvAddr := ":" + s.env.ApiPort
	server := http.Server{
		Addr:    srvAddr,
		Handler: s.mux,
	}
	return server.ListenAndServe()
}

func (s *HttpServer) RegisterRoutes() {
	routes.RegisterGeneralRoutes(s.mux, s.logger)
}
