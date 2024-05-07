package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env string `env-default:"local"`
	Database
	HTTPServer
}

type Database struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	Name     string `env:"POSTGRES_DB"`
	UserName string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:":8080"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"2s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

func LoadFromEnv() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse .env %v", err)
	}
	return &cfg, nil
}
