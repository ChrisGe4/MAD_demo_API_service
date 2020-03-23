package rest

import (
	"github.com/gin-gonic/gin"

	api "github.com/chrisge4/MAD_demo_API_service/pkg/api/rest"
	"github.com/chrisge4/MAD_demo_API_service/pkg/config"
)

func Routes(e *gin.Engine, sc *config.ServerConfig) {

	v1 := e.Group("/api/v1")
	v1.GET("/todo/:category", api.ListTodosFn(sc))
}
