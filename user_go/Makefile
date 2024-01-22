export TEST_POSTGRES_DB_URI=postgresql://user:secret@localhost:5432/postgres?sslmode=disable

TEST_POSTGRES_DB_URI=postgresql://user:secret@localhost:5432/postgres?sslmode=disable
CFG_PATH = "./configs/local.yml"
GO = go
GOBUILD = $(GO) build
GORUN = $(GO) run
GOTEST=$(GO) test
BINARY_NAME = app.exe
APP_PATH = ./cmd/main.go
TEST_POSTGRES_PATH = ./internal/storage/postgres
MIGRATION_PATH=./migrations/

all: build run

sql: sql_down sql_up

sql_down:
	goose -dir $(MIGRATION_PATH) postgres $(TEST_POSTGRES_DB_URI) reset

sql_up:
	goose -dir $(MIGRATION_PATH) postgres $(TEST_POSTGRES_DB_URI) up

test:
	$(GOTEST) $(TEST_POSTGRES_PATH)

build:
	$(GOBUILD) -o $(BINARY_NAME) $(APP_PATH)

run:
	$(BINARY_NAME) --config=$(CFG_PATH)

proto: proto_auth proto_perm

proto_auth:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/auth.proto

proto_perm:
	protoc --go_out=pb/ --go-grpc_out=pb/ protos/permissions.proto
