package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"server/service/api"
	"server/static"
)

const (
	Version      = "See documentation"
	Commit       = "See documentation"
	BuildDate    = "See documentation"
	isDebugBuild = "See documentation"
)

type Server struct {
	logger    *zap.Logger
	service   api.RestAPIService
	transport *echo.Echo
	jwtConfig *middleware.JWTConfig
}

func NewServer(l *zap.Logger, s api.RestAPIService, t *echo.Echo, jwtConfig *middleware.JWTConfig) Server {
	srv := Server{
		logger:    l,
		service:   s,
		transport: t,
		jwtConfig: jwtConfig,
	}
	srv.configureHandlers()
	return srv
}

func (s *Server) configureHandlers() {
	s.configureServiceHandlers()
}

func (s *Server) Serve(ctx context.Context, errs chan error) {
	s.logger.Info("Server", zap.String("name", static.ServerName), zap.String("version", Version), zap.String("commit", Commit), zap.String("build_date", BuildDate), zap.String("debug_build", isDebugBuild))
	port := ":5713"
	s.logger.Info("Start listening", zap.String("port", port))
	err := s.transport.Start(port)
	if err != nil {
		errsd := s.transport.Shutdown(ctx)
		if errsd != nil {
			s.logger.Error("shutting down transport", zap.Error(errsd))
		}
		errs <- err
	}
}
