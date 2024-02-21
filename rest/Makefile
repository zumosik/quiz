export CONFIG_PATH=./configs/config.yml

GO = go
GOBUILD = $(GO) build
GORUN = $(GO) run
GOTEST=$(GO) test
BINARY_NAME = app
APP_PATH = ./cmd/main.go

.DEFAULT_GOAL := all

all: build run

build:
	$(GOBUILD) -o $(BINARY_NAME) $(APP_PATH)

run:
	./$(BINARY_NAME)

proto: proto_files proto_auth

proto_files:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/files.proto

proto_auth:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/auth.proto


