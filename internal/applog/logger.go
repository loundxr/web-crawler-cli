package applog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func New(path string) (logger *slog.Logger, cleanup func(), err error) {
	dir := filepath.Dir(path)

	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, nil, fmt.Errorf("create log directory: %w", err)
		}
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(handler)

	cleanup = func() {
		file.Close()
	}

	return logger, cleanup, nil
}
