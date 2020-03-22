package rpc

import (
	"encoding/json"

	pb "github.com/chrisge4/MAD_demo_API_service/pkg/rpc/proto"
)

var mockData = []byte(`[{
    "description": "Attend meeting",
	"category": {
        "name": "work"
    },
	"id": 1
}, {
	"description": "Write code",
    "category": {
        "name": "work"
    },
	"id": 2
}]`)

type Server struct {
}

func (s *Server) ListTodos(c *pb.Category, stream pb.Todo_ListTodosServer) error {
	var todos []*pb.TodoItem
	err := json.Unmarshal(mockData, &todos)
	if err != nil {
		return err
	}

	for _, t := range todos {
		if t.Category.Name == c.Name {

			if err := stream.Send(t); err != nil {
				return err
			}
		}
	}

	return nil
}
