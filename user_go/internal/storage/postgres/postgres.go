package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"user_service/internal/domain/models"
	"user_service/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// MustOpenPostgresDB is helper func to get db, will throw panic when error
func MustOpenPostgresDB(dbURI string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		panic("cant open db" + err.Error())
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		panic("cant ping db " + err.Error())
	}

	return db
}

func (s *Storage) SaveUser(ctx context.Context, u models.User) (models.User, error) {
	if len(u.Password) < 3 || len(u.Email) < 3 {
		return models.User{}, storage.ErrEmptyFields
	}
	if err := u.EncryptPassword(); err != nil {
		return models.User{}, err
	}
	id := uuid.New().String()
	u.ID = id

	tx, err := s.db.Begin()
	// Make sure to close transaction if something goes wrong.
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := "INSERT INTO users(id, email, enc_password) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, query, u.ID, u.Email, u.EncPassword)
	if err != nil {
		return models.User{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (s *Storage) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := "SELECT id, enc_password FROM users WHERE email = $1"
	var id, encPassword string
	err := s.db.QueryRowContext(ctx, query, email).Scan(&id, &encPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrNotFound
		}
		return models.User{}, err
	}

	return models.User{
		ID:          id,
		Email:       email,
		EncPassword: []byte(encPassword),
	}, nil
}
