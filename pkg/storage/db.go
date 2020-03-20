package storage

import (
	"context"
	"io"
)

type Database interface {
	Get(ctx context.Context, uri string) (io.ReadCloser, error)
}
