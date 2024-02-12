package main

import (
	"context"
	"files/internal/config"
	"files/internal/grpc"
	"files/internal/grpc/files"
	"files/internal/storage/firebase_file_storage"
	"files/lib/slogpretty"
	"files/lib/utils"
	firebase "firebase.google.com/go"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	ctx := context.Background()

	firebaseApp, err := firebase.NewApp(ctx, cfg.StorageCfg, cfg.StorageOptions)
	if err != nil {
		panic(err)
	}
	client, err := firebaseApp.Storage(ctx)
	if err != nil {
		panic(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		panic(err)
	}

	grpcSrv := grpc.CreateGrpcServer(log)
	storage := firebase_file_storage.New(bucket)
	files.Register(grpcSrv, storage, log)

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
