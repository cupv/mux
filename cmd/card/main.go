package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cupv/mux/internal/config"
	cardHttp "github.com/cupv/mux/internal/delivery/http"
	"github.com/cupv/mux/internal/repository"
	"github.com/cupv/mux/internal/usecase"
	mysql "github.com/cupv/mux/pkg/mysql"
	"github.com/gorilla/mux"
)

// Runner defines the interface for running and shutting down a server
type Runner interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// Server is a concrete implementation of ServerRunner
type Server struct {
	server *http.Server
}

// NewRealServer creates a new RealServer instance
func NewRealServer(addr string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// setupLogger initializes the structured logger
func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// serveGracefully manages the server's lifecycle with graceful shutdown
func serveGracefully(runner Runner, logger *slog.Logger, addr string) int {
	// Channels for errors and shutdown signals
	errChan := make(chan error, 1)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("Starting Card Service", "address", addr)
		if err := runner.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for either an error or shutdown signal
	select {
	case err := <-errChan:
		logger.Error("Server failed to start", "error", err)
		return 1
	case <-shutdownChan:
		logger.Info("Shutdown signal received, initiating graceful shutdown")

		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Perform graceful shutdown
		if err := runner.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown failed", "error", err)
			return 1
		}
		logger.Info("Server gracefully stopped")
		return 0
	}
}

func main() {
	// Parse command-line flags for port
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// Load config
	config, err := config.LoadConfig()
	if err != nil {
		return
	}

	// Set up logger
	logger := setupLogger()
	slog.SetDefault(logger)

	// Set up db
	db, dbErr := mysql.Serve(config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	if dbErr != nil {
		return
	}
	defer db.Close()

	// Set up layers for clean arch
	repository := repository.NewCardRepository(db.Conn)
	service := usecase.NewCardUsecase(repository)
	handler := cardHttp.NewCardHandler(service)

	// Initialize router and server
	router := mux.NewRouter()
	router.HandleFunc("/cards", handler.GetCards).Methods("GET")
	router.HandleFunc("/card", handler.Create).Methods("POST")

	addr := ":" + *port
	server := NewRealServer(addr, router)

	// Run server and exit with appropriate code
	os.Exit(serveGracefully(server, logger, addr))
}
