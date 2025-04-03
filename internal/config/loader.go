package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "./.env"

func Load[T any]() (*T, error) {
	if _, err := os.Stat(envPath); err != nil {
		return nil, fmt.Errorf("%w: %w", errStat, err)
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("%w: %w", errLoad, err)
	}

	var cfg T
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", errRead, err)
	}

	return &cfg, nil
}
