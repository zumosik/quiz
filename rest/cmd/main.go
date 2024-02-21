package main

import (
	"fmt"
	"log/slog"
	"os"
	"rest_grpc/internal/config"
	"rest_grpc/internal/server"
	"rest_grpc/utils/log/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustGetConfig()
	log := setupLogger(cfg.Env)
	log.Debug("config readed, log configured")
	s := server.MustNew(log, cfg.Services.Files, cfg.Services.Users, cfg.Timeout)

	log.Info(fmt.Sprintf("Running server on port %s", cfg.Addr))
	err := s.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
