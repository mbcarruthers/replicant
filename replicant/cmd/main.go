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
	base := r.Group("/")
	api := base.Group("/api")
	v1 := api.Group("v1")

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "this is working",
		})
		return
	})

	v1.GET("/health", func(c *gin.Context) {
		if data, err := io.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":     err.Error(),
				"message": "first error",
			})
			return
		} else if _, err = c.Writer.Write(data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "second error",
			})
			return
		}
		return
	})

	if err := r.Run(port); err != nil {
		log.Fatalf(err.Error())
	}
}
