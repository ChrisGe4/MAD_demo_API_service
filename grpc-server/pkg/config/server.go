package config

import (
	"log"

	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/pb"
)

type ServerConfig struct {
	rpcClient pb.TodoClient
	debug     bool
}

func New(rpcClient pb.TodoClient) *ServerConfig {
	return &ServerConfig{rpcClient: rpcClient}
}

func (c *ServerConfig) IsDebug() bool {
	return c.debug
}

func (c *ServerConfig) SetDebug(debug bool) {
	c.debug = debug
}

func (c *ServerConfig) RpcClient() pb.TodoClient {
	if c.rpcClient == nil {
		log.Fatal("config: RPC Client not initialised")
	}

	return c.rpcClient
}
