package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"pingpongpoc/internal/constants"
	pingpongFilter "pingpongpoc/internal/domain/pingpong/filter"
	pingpongValidator "pingpongpoc/internal/domain/pingpong/validator"
	pingpongService "pingpongpoc/internal/service/pingpong"
	"pingpongpoc/internal/transport/http/handler"
)

// Server encapsulates the HTTP server over a Unix socket.
type Server struct {
	log        *slog.Logger
	server     *http.Server
	socketPath string
}

func NewServer(socketPath string, logger *slog.Logger) *Server {
	router := chi.NewRouter()

	service := pingpongService.NewService(
		pingpongFilter.New(logger),
		pingpongValidator.New(logger),
		logger,
	)
	pingHandler := handler.New(service, logger)

	router.Get("/", pingHandler.Ping)

	return &Server{
		server: &http.Server{ //nolint:exhaustruct
			Handler:           router,
			ReadHeaderTimeout: constants.ReadHeadTimeout,
		},
		log:        logger,
		socketPath: socketPath,
	}
}

func (s *Server) Start() error {
	if err := os.RemoveAll(s.socketPath); err != nil {
		s.log.Error("Failed to remove existing socket file", "error", err)

		return fmt.Errorf("failed to remove existing socket file: %w", err)
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		s.log.Error("Server failed to start", "error", err)

		return fmt.Errorf("server failed to start: %w", err)
	}
	defer listener.Close()

	s.log.Info("Starting server", "socketPath", s.socketPath)

	err = s.server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Error("Server failed during operation", "error", err)

		return fmt.Errorf("server failed during operation: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Failed shutting down server", "error", err)

		return fmt.Errorf("shutting down error: %w", err)
	}

	s.log.Info("HTTP Unix Server shut down gracefully")

	if err := os.Remove(s.socketPath); err != nil {
		s.log.Error("Failed to remove socket file after shutdown", "error", err)
	}

	return nil
}
