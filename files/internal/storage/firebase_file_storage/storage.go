package firebase_file_storage

import "cloud.google.com/go/storage"

type Storage struct {
	*storage.BucketHandle
}
