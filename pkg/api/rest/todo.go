package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chrisge4/MAD_demo_API_service/pkg/config"
	"github.com/chrisge4/MAD_demo_API_service/pkg/pb"
)

func ListFn(cfg *config.ServerConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		name := c.Param("category")
		category := &pb.Category{
			Name: name,
		}

		stream, err := cfg.RpcClient().List(ctx, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var result []*pb.Task
		for {
			t, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("%v.ListTodos(%v) = _, %v", cfg.RpcClient(), category, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result = append(result, t)
		}
		c.JSON(http.StatusOK, result)
	}
}

func AddFn(cfg *config.ServerConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		log.Println("Adding a new todo")
		var t pb.Task
		err := c.BindJSON(&t)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		r, err := cfg.RpcClient().Add(ctx, &t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, r)
	}
}
