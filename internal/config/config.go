package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mkz-Prog/yadro-telecom-test/internal/domain"
)

// Load читает и парсит конфигурационный файл по указанному пути.
func Load(path string) (*domain.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var cfg domain.Config

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config json: %w", err)
	}

	return &cfg, nil
}
