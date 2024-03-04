package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	pingpongFilter "pingpongpoc/internal/domain/pingpong/filter"
	pingpongValidator "pingpongpoc/internal/domain/pingpong/validator"
	pingpongService "pingpongpoc/internal/service/pingpong"
	"pingpongpoc/internal/transport/unix/handler"
)

// ! IMPORTANT !
// IS NOT FULLY WORKED CODE
// Graceful shutdown doesn't work properly

// Server represents a server listening on a Unix socket.
type (
	Handler interface {
		Ping(ctx context.Context, conn net.Conn)
	}

	//nolint:govet
	Server struct {
		log            *slog.Logger
		handler        Handler
		listener       net.Listener
		shutdownCtx    context.Context //nolint:containedctx // Context for lifecycle the server
		shutdownCancel context.CancelFunc
		connections    map[net.Conn]struct{}
		connLock       sync.Mutex
		socketPath     string
	}
)

func NewServer(socketPath string, logger *slog.Logger) *Server {
	service := pingpongService.NewService(
		pingpongFilter.New(logger),
		pingpongValidator.New(logger),
		logger,
	)
	pingHandler := handler.New(service, logger)

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	return &Server{
		log:            logger,
		handler:        pingHandler,
		listener:       nil,
		connections:    make(map[net.Conn]struct{}),
		connLock:       sync.Mutex{},
		socketPath:     socketPath,
		shutdownCtx:    shutdownCtx,
		shutdownCancel: shutdownCancel,
	}
}

func (s *Server) Start() error {
	err := os.RemoveAll(s.socketPath)
	if err != nil {
		s.log.Error("Failed to remove existing socket file", "error", err)

		return fmt.Errorf("failed to remove existing socket file: %w", err)
	}

	s.listener, err = net.Listen("unix", s.socketPath)
	if err != nil {
		s.log.Error("Server failed to start", "error", err)

		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	defer s.listener.Close()

	s.log.Info("Starting server", "socketPath", s.socketPath)

	go func() {
		<-s.shutdownCtx.Done()
		s.listener.Close() // Close the listener on shutdown signal
	}()

	return s.Listen()
}

// TODO await all (maybe review this shutdown system“)
func (s *Server) Shutdown() {
	s.log.Debug("Shutdown")
	s.shutdownCancel() // Signal all operations to shutdown
}

func (s *Server) Listen() error {
	for {
		conn, listenErr := s.listener.Accept()
		if listenErr != nil {
			select {
			case <-s.shutdownCtx.Done():
				s.log.Info("Shuting down server successfully")

				return nil // Exit if shutdown is initiated
			default:
				s.log.Error("Listening", "error", listenErr)

				continue
			}
		}

		// TODO make better connections registration

		s.connLock.Lock()
		s.connections[conn] = struct{}{}
		s.connLock.Unlock()

		// TODO add routing

		go func(connection net.Conn) {
			connectionContext, cancel := context.WithCancel(s.shutdownCtx)
			// TODO make correct graceful shutdown
			defer func() {
				cancel()
				connection.Close()

				// Cleaning active connection from pull
				s.connLock.Lock()
				delete(s.connections, connection)
				s.connLock.Unlock()
			}()
			s.handler.Ping(connectionContext, connection)
		}(conn)
	}
}
