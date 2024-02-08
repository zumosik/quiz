package main

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"io"
	"os"
)

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("opt.json")
	config := &firebase.Config{
		DatabaseURL:   "gs://files-saver-2233.appspot.com",
		StorageBucket: "files-saver-2233.appspot.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		panic(err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		panic(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		panic(err)
	}

	f, err := os.Open("download.png")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	obj := bucket.Object(f.Name())
	w := obj.NewWriter(ctx)
	id := uuid.New()
	w.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer func(w *storage.Writer) {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}(w)

	if _, err := io.Copy(w, f); err != nil {
		panic(err)
	}

	fmt.Print("success")
}
