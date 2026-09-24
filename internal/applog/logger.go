package applog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/loundxr/web-crawler-cli/internal/applog/handlers"
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

	formatHandler := handlers.NewFormatHandler(file)
	logger = slog.New(formatHandler)

	cleanup = func() {
		file.Close()
	}

	return logger, cleanup, nil
}
