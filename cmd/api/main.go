package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chrisge4/MAD_demo_API_service/pkg/config"
	"github.com/chrisge4/MAD_demo_API_service/pkg/http/rest"
	storage "github.com/chrisge4/MAD_demo_API_service/pkg/storage/nosql"
)

func main() {
	ctx := context.Background()
	db, err := storage.NewGcs(ctx, "", "gcore")
	fatal(err)
	cfg := config.New(db)
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
