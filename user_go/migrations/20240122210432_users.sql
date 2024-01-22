-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
  id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE ,
  enc_password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
