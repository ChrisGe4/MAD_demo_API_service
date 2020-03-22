package config

import (
	"log"

	todo "github.com/chrisge4/MAD_demo_API_service/pkg/rpc/proto"
)

type ServerConfig struct {
	rpcClient todo.TodoClient
	debug     bool
}

func New(rpcClient todo.TodoClient) *ServerConfig {
	return &ServerConfig{rpcClient: rpcClient}
}

func (c *ServerConfig) IsDebug() bool {
	return c.debug
}

func (c *ServerConfig) SetDebug(debug bool) {
	c.debug = debug
}

func (c *ServerConfig) RpcClient() todo.TodoClient {
	if c.rpcClient == nil {
		log.Fatal("config: RPC Client not initialised")
	}

	return c.rpcClient
}
