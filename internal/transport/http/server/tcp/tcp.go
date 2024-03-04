package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"pingpongpoc/internal/constants"
	pingpongFilter "pingpongpoc/internal/domain/pingpong/filter"
	pingpongValidator "pingpongpoc/internal/domain/pingpong/validator"
	pingpongService "pingpongpoc/internal/service/pingpong"
	"pingpongpoc/internal/transport/http/handler"
)

// Server encapsulates the HTTP server.
type Server struct {
	log    *slog.Logger
	server *http.Server
}

// NewServer creates a new Server instance.
func NewServer(addr string, logger *slog.Logger) *Server {
	router := chi.NewRouter()

	service := pingpongService.NewService(
		pingpongFilter.New(logger),
		pingpongValidator.New(logger),
		logger,
	)
	pingHandler := handler.New(service, logger)

	router.Get("/", pingHandler.Ping)

	//nolint:exhaustruct
	return &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: constants.ReadHeadTimeout,
		},
		log: logger,
	}
}

// Start tarts the server.
func (s *Server) Start() error {
	s.log.Info("Starting server", "address", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Error("Server failed to start", "error", err)

		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error(
			"Failed shutting down server",
			"error", err,
		)

		return fmt.Errorf("shutting down error: %w", err)
	}

	s.log.Info("HTTP TCP IP Server shut down gracefully")

	return nil
}
