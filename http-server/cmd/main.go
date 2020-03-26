package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/chrisge4/MAD_demo_API_service/http-server/pkg/config"
	"github.com/chrisge4/MAD_demo_API_service/http-server/pkg/pb"
	"github.com/chrisge4/MAD_demo_API_service/http-server/pkg/routes"
)

var (
	host string
	port string
)

const (
	defGrpcServerAddr = "localhost:8090"
	GrpcHostEnv       = "GRPC_SERVER_HOST"
	GrpcPortEnv       = "GRPC_SERVER_PORT"
)

func main() {
	flag.StringVar(&host, "host", "", "host of grpc server")
	flag.StringVar(&port, "port", "", "port of grpc server")
	flag.Parse()
	if host == "" {
		if host, port = os.Getenv(GrpcHostEnv), os.Getenv(GrpcPortEnv); host == "" || port == "" {
			log.Printf("host %v port %v \n", host, port)
			host = defGrpcServerAddr
			log.Printf("addr %v \n", host)

		} else {
			host = strings.Join([]string{host, port}, ":")
		}
	}
	log.Printf("GRPC server address is %v \n", host)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewTodoClient(conn)
	cfg := config.New(client)
	cfg.SetDebug(true)
	s := gin.Default()
	s.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Rest server running.")
	})
	routes.Routes(s, cfg)

	//http.ListenAndServe("8080", gin)
	s.Run(":8081")
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
