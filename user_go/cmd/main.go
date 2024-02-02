package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user_service/internal/config"
	"user_service/internal/grpc"
	"user_service/internal/grpc/auth"
	"user_service/internal/storage/postgres"
	"user_service/lib/slogpretty"
	"user_service/lib/utils"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Config read!")

	grpcSrv := grpc.CreateGrpcServer(log)

	db := postgres.MustOpenPostgresDB(cfg.PostgresStorageURI)
	storage := postgres.New(db)

	log.Info("Storage created!")

	auth.Register(grpcSrv, storage, log, cfg.TokenSecret, cfg.TokenTTL)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Error("cant listen (tcp)", utils.WrapErr(err))
	}

	log.Info(fmt.Sprintf("Server started on port %d", cfg.GRPC.Port))

	go func() {
		if err := grpcSrv.Serve(l); err != nil {
			log.Error("cant serve server", utils.WrapErr(err))
		}
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	grpcSrv.GracefulStop()

	log.Info("Server stopped...")
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
