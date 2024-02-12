package models

import "time"

type File struct {
	ID        string
	Name      string
	CreatedAt time.Time
	Bytes     []byte
}
