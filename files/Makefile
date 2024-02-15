export FIREBASE_CFG_PATH=configs/firebase_opt.json
export CFG_PATH=configs/local.yml

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

proto: proto_files

proto_files:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/files.proto
