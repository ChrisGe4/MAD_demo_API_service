package api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chrisge4/MAD_demo_API_service/pkg/config"
)

func GetTodoFn(cfg *config.ServerConfig) func(*gin.Context) {
	return func(c *gin.Context) {

		ctx := context.Background()

		id := c.Param("id")
		rc, err := cfg.Db().Get(ctx, id)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		defer rc.Close()

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Writer.Write(data)

	}
}
