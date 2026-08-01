// Package main is the entry point for the dumaVote application.
package main

import (
	"bytes"
	"context"
	"dumaVote/server"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

// envNames lists required environment variables.
var envNames = []string{"APP_API_KEY", "PERSONAL_API_KEY"}

func checkEnv() {
	for _, name := range envNames {
		if os.Getenv(name) == "" {
			log.Fatalf("Required environment variable %s is not set", name)
		}
	}
}

// StdoutHandler is a custom slog.Handler that writes structured logs to stdout.
type StdoutHandler struct {
	opts *slog.HandlerOptions
	mu   sync.Mutex
	out  io.Writer
}

// NewStdoutHandler creates a new StdoutHandler.
func NewStdoutHandler(w io.Writer, opts *slog.HandlerOptions) *StdoutHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &StdoutHandler{opts: opts, out: w}
}

// Enabled reports whether the handler handles records at the given level.
func (h *StdoutHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// WithAttrs returns a new Handler with the given attributes added.
func (h *StdoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Simplified: attributes are appended in Handle.
	return h
}

// WithGroup returns a new Handler with the given group prepended to future group names.
func (h *StdoutHandler) WithGroup(name string) slog.Handler {
	// Simplified: group is ignored.
	return h
}

// Handle processes a log record.
func (h *StdoutHandler) Handle(_ context.Context, r slog.Record) error {
	buf := &bytes.Buffer{}

	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := strings.ToUpper(r.Level.String())

	var sourceStr string
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		fileName := filepath.Base(f.File)
		sourceStr = fmt.Sprintf("%s:%d", fileName, f.Line)
	}

	fmt.Fprintf(buf, "[%s][%s][%s] %s", timeStr, levelStr, sourceStr, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(buf, " %s=%v", a.Key, a.Value.Any())
		return true
	})

	fmt.Fprintln(buf)

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(buf.Bytes())
	return err
}

func main() {
	checkEnv()

	handler := NewStdoutHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(handler)

	port := "8080"

	srv := server.NewDumaVotesServer(
		os.Getenv("APP_API_KEY"),
		os.Getenv("PERSONAL_API_KEY"),
		logger,
	)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.HandlerFunc(srv.MainHandler),
	}

	// Graceful shutdown on SIGINT/SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info(fmt.Sprintf("Start server on http://localhost:%s", port))

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	if err := srv.Close(); err != nil {
		logger.Error("Database close error", "error", err)
	}

	logger.Info("Server exited")
}
