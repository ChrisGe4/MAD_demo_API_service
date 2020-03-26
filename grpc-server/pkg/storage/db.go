package storage

import (
	"context"

	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/pb"
)

type Database interface {
	//Get(ctx context.Context, uri string) (io.ReadCloser, error)
	Add(context.Context, *pb.Task) error
	List(context.Context, *pb.Category) ([]*pb.Task, error)
	Close() error
}
