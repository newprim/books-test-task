package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}

	HTTP struct {
		Port     string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		MaxPRD   int           `env-required:"true" yaml:"max_rpd" env:"HTTP_MAX_RPD"`
		Duration time.Duration `env-required:"true" yaml:"duration" env:"HTTP_DURATION"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
)

func NewConfig() (Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
