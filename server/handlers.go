package server

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (s *Server) configureServiceHandlers() {
	s.transport.GET("/ping", s.Ping)
}

func (s *Server) Ping(c echo.Context) error {
	type result struct {
		Result string
	}

	res := &result{Result: "pong!"}
	s.logger.Info("Heart beat ping!", zap.String("path", c.Path()))
	return c.JSON(http.StatusOK, res)
}