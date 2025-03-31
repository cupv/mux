package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "os"
    "syscall"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockRunner struct {
    mock.Mock
}

func (m *MockRunner) ListenAndServe() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockRunner) Shutdown(ctx context.Context) error {
    args := m.Called(ctx)
    return args.Error(0)
}

func TestServerGracefulShutdown(t *testing.T) {
    mockRunner := new(MockRunner)
    logger := setupLogger()
    addr := ":8080"

    mockRunner.On("ListenAndServe").Return(http.ErrServerClosed).Once()
    mockRunner.On("Shutdown", mock.Anything).Return(nil).Once()

    go func() {
        time.Sleep(1 * time.Second)
        process, _ := os.FindProcess(os.Getpid())
        process.Signal(syscall.SIGINT)
    }()

    exitCode := serveGracefully(mockRunner, logger, addr)
    assert.Equal(t, 0, exitCode, "Expected graceful shutdown to return 0 exit code")
    mockRunner.AssertExpectations(t)
}

func TestServerListenAndServeError(t *testing.T) {
    mockRunner := new(MockRunner)
    logger := setupLogger()
    addr := ":8080"

    // Simulate server failing to start
    mockRunner.On("ListenAndServe").Return(assert.AnError).Once()

    exitCode := serveGracefully(mockRunner, logger, addr)
    assert.Equal(t, 1, exitCode, "Expected failure exit code 1")
    mockRunner.AssertExpectations(t)
}

func TestNewRealServer(t *testing.T) {
    handler := http.NewServeMux()
    server := NewRealServer(":8080", handler)

    assert.NotNil(t, server, "Server should not be nil")
    assert.Equal(t, ":8080", server.server.Addr, "Server address should match")
}

func TestServerShutdown(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    defer ts.Close()

    server := NewRealServer(ts.Listener.Addr().String(), ts.Config.Handler)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := server.Shutdown(ctx)
    assert.NoError(t, err, "Expected no error on shutdown")
}