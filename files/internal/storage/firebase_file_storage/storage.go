package firebase_file_storage

import (
	"context"
	"files/internal/domain/models"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

const (
	IDMetadataName        = "id"
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

	obj := s.bucket.Object(f.Name)
	w := obj.NewWriter(ctx)
	w.ObjectAttrs.Metadata = map[string]string{"id": id, "data": "lol"}
	defer func(w *storage.Writer) {
		_ = w.Close()
	}(w)

	_, err := w.Write(f.Bytes)
	if err != nil {
		return models.File{}, err
	}

	//TODO: write this
}
