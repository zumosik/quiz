package storage

import "errors"

var (
	ErrEmptyFields = errors.New("empty fields")
	ErrNotFound    = errors.New("not found")
)
