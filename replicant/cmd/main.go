package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	api := r.Group("/api/v1")

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "this is working",
		})
		return
	})
	
	if err := r.Run(port); err != nil {
		log.Fatalf(err.Error())
	}
}
