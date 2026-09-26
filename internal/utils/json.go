package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/loundxr/web-crawler-cli/internal/model"
)

func WriteResults(nodes []*model.Node, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create results directory: %w", err)
	}

	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal nodes: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
