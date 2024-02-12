package firebase_file_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"files/internal/domain/models"
	"github.com/google/uuid"
)

const (
	IDMetadataName        = "id"
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

	obj := s.bucket.Object(id + f.UserID + f.Name)
	w := obj.NewWriter(ctx)
	w.ObjectAttrs.Metadata = map[string]string{
		IDMetadataName:        id,
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
