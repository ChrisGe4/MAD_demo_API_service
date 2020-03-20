package config

import (
	"github.com/chrisge4/MAD_demo_API_service/pkg/storage"
)

type ServerConfig struct {
	db    storage.Database
	debug bool
}

func New(db storage.Database) *ServerConfig {
	return &ServerConfig{db: db}
}

func (c *ServerConfig) IsDebug() bool {
	return c.debug
}

func (c *ServerConfig) SetDebug(debug bool) {
	c.debug = debug
}
