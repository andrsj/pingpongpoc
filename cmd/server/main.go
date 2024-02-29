package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	httpServer "pingpongpoc/internal/transport/http/server"
	unixServer "pingpongpoc/internal/transport/unix/server"
)

const (
	shutdownTimeout            = 5 * time.Second
	httpListenAddressTCPServer = ":8080"
	pathToSocketHTTPUnixServer = "/tmp/pingpong1.sock"
	pathToSocketUnixServer     = "/tmp/pingpong2.sock"
)

//nolint:funlen,gocognit
func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}

	// JSON or Text handler
	// logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	server1 := httpServer.NewTCPServer(httpListenAddressTCPServer, logger)
	server2 := httpServer.NewUnixServer(pathToSocketHTTPUnixServer, logger)
	server3 := unixServer.New(pathToSocketUnixServer, logger)

	mainContext, mainCancelFunc := context.WithCancel(context.Background())
	defer mainCancelFunc()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info("Received shutdown signal")

		// Create a context for the shutdown process with a timeout.
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelShutdown()

		waitGroup := sync.WaitGroup{}

		waitGroup.Add(3) //nolint:gomnd

		go func() {
			defer waitGroup.Done()

			if err := server1.Shutdown(shutdownCtx); err != nil {
				logger.Error("Error during server1 shutdown", "error", err)
			} else {
				logger.Info("Server1 shut down gracefully")
			}
		}()

		go func() {
			defer waitGroup.Done()

			if err := server2.Shutdown(shutdownCtx); err != nil {
				logger.Error("Error during server2 shutdown", "error", err)
			} else {
				logger.Info("Server2 shut down gracefully")
			}
		}()

		go func() {
			defer waitGroup.Done()

			server3.Shutdown()
			logger.Info("Server3 shut down gracefully")
		}()

		waitGroup.Wait()
		mainCancelFunc()
	}()

	go func() {
		if err := server1.Start(); err != nil {
			logger.Error("Error starting HTTP/TCP-IP server", "error", err)
		}
	}()

	go func() {
		if err := server2.Start(); err != nil {
			logger.Error("Error starting HTTP/Unix server1", "error", err)
		}
	}()

	go func() {
		if err := server3.Start(); err != nil {
			logger.Error("Error starting Unix server2", "error", err)
		}
	}()

	<-mainContext.Done()
}
