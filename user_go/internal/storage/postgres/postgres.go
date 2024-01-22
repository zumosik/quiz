package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"user_service/internal/domain/models"
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
	// TODO: add empty email validation
	if err := u.EncryptPassword(); err != nil {
		return models.User{}, err
	}
	id := uuid.New().String()
	u.ID = id

	query := "INSERT INTO users(id, email, enc_password) VALUES ($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, u.ID, u.Email, u.EncPassword)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
