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
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func isExist(env_name string) bool {
	_, isExists := os.LookupEnv(env_name)
	return isExists
}

func checkEnv() {
	if !isExist("APP_API_KEY") {
		log.Fatal("App api key is not set. APP_API_KEY env is empty")
	}
	if !isExist("PERSONAL_API_KEY") {
		log.Fatal("Personal api key is not set. PERSONAL_API_KEY env is empty")
	}
}

// CustomHandler реализует интерфейс slog.Handler с кастомным форматом
type CustomHandler struct {
	opts *slog.HandlerOptions
	mu   sync.Mutex
	out  io.Writer
}

func NewCustomHandler(w io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &CustomHandler{
		opts: opts,
		out:  w,
	}
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Для простоты игнорируем дополнительные атрибуты
	// В production можно хранить их и выводить
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	buf := &bytes.Buffer{}

	// Формат времени: 2026-08-02 15:30:45
	timeStr := r.Time.Format("2006-01-02 15:04:05")

	// Уровень логирования в верхнем регистре
	levelStr := strings.ToUpper(r.Level.String())

	// Получаем информацию о файле и строке
	var sourceStr string
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		// Только имя файла без пути
		fileName := filepath.Base(f.File)
		sourceStr = fmt.Sprintf("%s:%d", fileName, f.Line)
	}

	// Формируем сообщение: [Время][файл:строка] Сообщение
	fmt.Fprintf(buf, "[%s][%s][%s] %s", timeStr, levelStr, sourceStr, r.Message)

	// Добавляем атрибуты если есть
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
	handler := NewCustomHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(handler)

	port := "8080"

	dumaVotesServer := server.NewDumaVotesServer(os.Getenv("APP_API_KEY"), os.Getenv("PERSONAL_API_KEY"), logger)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.HandlerFunc(dumaVotesServer.MainHandler),
	}
	logger.Info(fmt.Sprintf("Start server on http://localhost:%s", port))
	err := server.ListenAndServe()
	if err != nil {
		logger.Info("Http serve shut down.", "Error", err.Error())
	}
}
