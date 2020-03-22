package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chrisge4/MAD_demo_API_service/pkg/config"
	pb "github.com/chrisge4/MAD_demo_API_service/pkg/rpc/proto"
)

func GetTodoFn(cfg *config.ServerConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		name := c.Param("category")
		//rc, err := cfg.Db().Get(ctx, id)
		//if err != nil {
		//	c.AbortWithError(http.StatusNotFound, err)
		//	return
		//}
		//defer rc.Close()
		category := &pb.Category{
			Name: name,
		}
		stream, err := cfg.RpcClient().ListTodos(ctx, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var result []*pb.TodoItem
		for {

			t, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.ListTodos(%v) = _, %v", cfg.RpcClient(), category, err)
			}
			result = append(result, t)
		}

		c.JSON(http.StatusOK, result)

	}

}
