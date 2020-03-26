package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	rpc "github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/grpc"
	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/pb"
	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/storage"
)

const (
	defGcpProject = "cloud-build-delivery-mad"
	defCollection = "todo"
	GcpProjectEnv = "GCP_PROJECT_ID"
	CollectionEnv = "FIRESTORE_COLLECTION"
)

func main() {

	project := os.Getenv(GcpProjectEnv)
	collection := os.Getenv(CollectionEnv)
	if project == "" {
		project = defGcpProject
	}
	if collection == "" {
		collection = defCollection
	}
	gs := grpc.NewServer()
	db, err := storage.NewFirestore(context.Background(), "", project, collection)
	defer db.Close()
	if err != nil {
		log.Fatalf("failed to connect to firestore: %v", err)
	}
	ts := rpc.NewServer(db)
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
