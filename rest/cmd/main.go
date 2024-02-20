package main

import (
	"log/slog"
	"rest_grpc/internal/server"
)

func main() {
	s := server.MustNew(slog.Default(), "", "")

	err := s.Run(":3333")
	if err != nil {
		panic(err)
	}
}
