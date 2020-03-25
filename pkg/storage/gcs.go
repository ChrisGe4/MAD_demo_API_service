package storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/martian/log"
	"google.golang.org/api/option"
)

type Gcs struct {
	bucket *storage.BucketHandle
}

func NewGcs(ctx context.Context, credentialsFile, bucket string) (*Gcs, error) {
	var client *storage.Client
	var err error
	if credentialsFile == "" {
		client, err = storage.NewClient(ctx)
	} else {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))

	}

	if err != nil {
		return nil, err
	}

	return &Gcs{bucket: client.Bucket(bucket)}, nil
}

func (g *Gcs) Get(ctx context.Context, uri string) (io.ReadCloser, error) {
	log.Infof("%v", g.bucket)
	rc, err := g.bucket.Object(uri).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return rc, nil

}
