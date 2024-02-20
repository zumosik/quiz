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


