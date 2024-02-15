package firebase_file_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"files/internal/domain/models"
	storage2 "files/internal/storage"
	"github.com/google/uuid"
	"io"
	"time"
)

const (
	UserIDMetadataName    = "user_id"
	NameMetadataName      = "name"
	CreatedAtMetadataName = "created_at"
)

type Storage struct {
	bucket *storage.BucketHandle
}

func New(bucket *storage.BucketHandle) *Storage {
	return &Storage{
		bucket: bucket,
	}
}

func (s *Storage) UploadFile(ctx context.Context, f models.File) (models.File, error) {
	id := uuid.New().String()

	obj := s.bucket.Object(id)
	w := obj.NewWriter(ctx)
	w.ObjectAttrs.Metadata = map[string]string{
		UserIDMetadataName:    f.UserID,
		NameMetadataName:      f.Name,
		CreatedAtMetadataName: f.CreatedAt.String()}
	defer func(w *storage.Writer) {
		_ = w.Close()
	}(w)

	_, err := w.Write(f.Bytes)
	if err != nil {
		return models.File{}, err
	}

	return models.File{
		ID:        id,
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		Bytes:     f.Bytes,
	}, nil
}

func (s *Storage) GetFileById(ctx context.Context, id string) (models.File, error) {
	obj := s.bucket.Object(id)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return models.File{}, storage2.ErrNotFound
		}
		return models.File{}, nil
	}
	defer reader.Close()

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return models.File{}, storage2.ErrNotFound
		}
		return models.File{}, nil
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return models.File{}, err
	}

	// not parsing error because this error isn't critical
	parsedTime, _ := time.Parse(time.RFC3339Nano, attrs.Metadata[NameMetadataName])

	return models.File{
		ID:        id,
		UserID:    attrs.Metadata[UserIDMetadataName],
		Name:      attrs.Metadata[NameMetadataName],
		CreatedAt: parsedTime,
		Bytes:     bytes,
	}, nil

}
