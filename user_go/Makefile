CFG_PATH = "./configs/local.yml"
GO = go
GOBUILD = $(GO) build
GORUN = $(GO) run
BINARY_NAME = app.exe
APP_PATH = ./cmd/main.go

all: build run

build:
	$(GOBUILD) -o $(BINARY_NAME) $(APP_PATH)

run:
	$(BINARY_NAME) --config=$(CFG_PATH)

proto: proto_auth proto_perm

proto_auth:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/auth.proto

proto_perm:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/permissions.proto
