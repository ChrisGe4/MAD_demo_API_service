package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/chrisge4/MAD_demo_API_service/pkg/rpc"
	pb "github.com/chrisge4/MAD_demo_API_service/pkg/rpc/proto"
)

func main() {

	gs := grpc.NewServer()
	ts := &rpc.Server{}

	pb.RegisterTodoServer(gs, ts)
	//reflection.Register(gs)

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Starting server at port :8090")

	if err := gs.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
