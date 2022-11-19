package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

const (
	port = ":8000"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-type"},
		ExposeHeaders:    []string{"Content-Length", "Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api := r.Group("/api")

	api.GET("/echo", func(c *gin.Context) {
		if data, err := io.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		} else if _, err = c.Writer.Write(data); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	})

	if err := r.Run(port); err != nil {
		log.Fatalf(err.Error())
	}
}
