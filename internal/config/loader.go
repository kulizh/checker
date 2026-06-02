package config

import (
	"encoding/json"
	"os"

	"checker/internal/domain"
)

type Config map[string]domain.DomainConfig

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(b, &cfg)
	return cfg, err
}
