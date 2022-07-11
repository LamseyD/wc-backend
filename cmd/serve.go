package cmd

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"server/server"
	"server/service/api"
	"syscall"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Run in daemon mode",
	Long:    "Run as a microservice with HTTP transport",
	Aliases: []string{"run"},
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve() {
	zapLogger, _ := zap.NewDevelopment()
	zapLogger.Info("New development logger initiated")
	zapLogger.Info("Running wc-backend in daemon mode")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set blocking channel to keep server alive, exit if unblocked
	errChan := make(chan error)
	defer close(errChan)

	// Listen to interruption signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-termChan
		zapLogger.Info("got interruption signal", zap.String("signal", s.String()))
		zapLogger.Info("will force transport interruption by simulating an error")
		errChan <- errors.New("Signal termination")
	}()

	// create DB implementation --- all repositories go here

	s := api.RestAPIService{
		Logger: zapLogger,
	}

	// JWT middleware configuration if needed here
	//jwtMiddlewareConfig := configureJWTEchoMiddleware(zapLogger)
	jwtMiddlewareConfig := configureJWTEchoMiddleware(zapLogger)
	echoTransport := echo.New()
	svr := server.NewServer(zapLogger, s, echoTransport, jwtMiddlewareConfig)

	go svr.Serve(ctx, errChan)

	zapLogger.Info("service.go terminated from Echo error", zap.Error(<-errChan))
}

func configureJWTEchoMiddleware(logger *zap.Logger) *middleware.JWTConfig {
	// TODO: To be implemented
	return nil
}
