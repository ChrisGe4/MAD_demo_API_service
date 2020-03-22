package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/chrisge4/MAD_demo_API_service/pkg/config"
	"github.com/chrisge4/MAD_demo_API_service/pkg/http/rest"
	pb "github.com/chrisge4/MAD_demo_API_service/pkg/rpc/proto"
)

func main() {

	addr := flag.String("addr", "localhost:8082", "address of grpc server")
	flag.Parse()
	//ctx := context.Background()
	//db, err := storage.NewGcs(ctx, "", "gcore")
	//fatal(err)
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewTodoClient(conn)
	cfg := config.New(client)
	cfg.SetDebug(true)
	s := gin.Default()
	s.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Backend running.")
	})
	rest.Routes(s, cfg)

	//http.ListenAndServe("8080", gin)
	s.Run(":8081")
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
