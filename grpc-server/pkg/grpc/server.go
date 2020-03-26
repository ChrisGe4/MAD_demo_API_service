package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/pb"
	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/storage"
)

type Server struct {
	db storage.Database
}

func NewServer(db storage.Database) *Server {
	return &Server{db: db}
}
func (s *Server) List(c *pb.Category, stream pb.Todo_ListServer) error {
	grpclog.Infof("List todos in category %v", c.Name)
	tasks, err := s.db.List(context.Background(), c)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if err := stream.Send(t); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Add(ctx context.Context, t *pb.Task) (*pb.Task, error) {
	grpclog.Infof("Adding a new todo %v", t)
	t.Id = uuid.New().String()
	if err := s.db.Add(ctx, t); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return t, nil
}

func (s *Server) Close() error {
	return s.db.Close()
}
