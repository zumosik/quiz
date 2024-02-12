package models

import "time"

type File struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
	Bytes     []byte
}
