package main

import (
	"goxy/internal/config"
	"goxy/internal/proxy"
	useragent "goxy/internal/user-agent"
	"net/http"
	"os"

	"go.uber.org/zap"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.MustLoad()
	if err != nil {
		os.Exit(1)
	}

	// Конфигурируем Логгер
	zl := setupLogger(cfg.Env)
	defer zl.Sync()
	zl.Info("logger is configured")

	// Рандомими User-Agent
	userAgent, err := useragent.Random(zl)
	zl.Info("A random user-agent was received.", zap.String("User-Agent", userAgent.UserAgent))

	http.HandleFunc("/", proxy.Handle(zl, cfg, userAgent.UserAgent))
	zl.Info("Starting proxy server.", zap.String("Port:", cfg.ListenPort))
	http.ListenAndServe(":"+cfg.ListenPort, nil)
}

func setupLogger(env string) *zap.Logger {
	var zl *zap.Logger

	switch env {
	case envLocal:
		zl = zap.Must(zap.NewDevelopment())
	case envDev:
		zl = zap.Must(zap.NewDevelopment())
	case envProd:
		zl = zap.Must(zap.NewProduction())
	}

	return zl
}
