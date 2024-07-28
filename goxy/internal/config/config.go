package config

import (
	"errors"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Переменные окружения
type Config struct {
	Env            string        `env:"ENV"`
	Host           string        `env:"HOST"`
	ListenPort     string        `env:"LISTENPORT"`
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT"`
	TargetURL      string        `env:"TARGET_URL"`
	BackupURL      string        `env:"BACKUP_URL"`
	SecretKey      string        `env:"SECRET_KEY"`
}

func MustLoad() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Printf("couldn't parse the .env file: %v", err)
		return nil, err
	}

	if cfg.Env == "" {
		err := errors.New("Env is not set")
		log.Printf("The necessary variables are missing in the environment: %v", err)
		return nil, err
	}

	log.Print("config is loaded")
	return &cfg, nil
}
