package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	__port = 8000
	port   = fmt.Sprintf(":%d", __port)
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

	server := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen err:%s\n", err.Error())
		}
	}()
	quit := make(chan os.Signal)
	// Kill - (no param) - default send syscan11.SIGTERM
	// Kill - 2 = syscall.SIGINT
	// Kill - 9 = syscall.SIGKILL cannot catch with select statement
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Printf("Shutdown server\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown:%s \n", err.Error())
	}

	// catching ctx.Done() after 5 second timeout
	select {
	case <-ctx.Done():
		log.Printf("Timed out...\n")
	}
	log.Println("Server exiting...")
}
