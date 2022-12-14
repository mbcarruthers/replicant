package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DataStore struct {
	Conn *pgx.Conn
}

func NewDataStore(dbUrl string) (*DataStore, error) {
	if conn, err := pgx.Connect(context.Background(), dbUrl); err != nil {
		return nil, err
	} else {
		return &DataStore{
			Conn: conn,
		}, nil
	}
}

func (d *DataStore) CreateTestDatabase(ctx context.Context) error {
	tx, err := d.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	databaseOptions := map[string]string{
		"create_database": "CREATE DATABASE test_database",
		"create_table": "CREATE TABLE test_database.test_table(" +
			"id uuid NOT NULL DEFAULT gen_random_uuid()," +
			"name STRING NOT NULL)",
	}
	if _, err := tx.Exec(ctx, databaseOptions["create_database"]); err != nil {
		return err
	} else if _, err = tx.Exec(ctx, databaseOptions["create_table"]); err != nil {
		return err
	} else {
		return tx.Commit(ctx)
	}
}

func (d *DataStore) InsertTestElement(ctx context.Context) error {
	tx, err := d.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	insertionString := fmt.Sprintf("INSERT INTO test_database.test_table(name) VALUES($1)")
	if t, err := tx.Exec(ctx, insertionString, "Decker"); err != nil {
		return err
	} else if t.RowsAffected() == 0 {
		return fmt.Errorf("Did not affect any rows")
	} else {
		return tx.Commit(ctx)
	}
}

func (d *DataStore) QueryTestElement(ctx context.Context) (string, error) {
	queryString := fmt.Sprintf("SELECT name FROM test_database.test_table WHERE name = $1")
	if rows, err := d.Conn.Query(ctx, queryString, "Decker"); err != nil {
		return "", err
	} else {
		var names []string
		for rows.Next() {
			var name string
			if err = rows.Scan(&name); err != nil {
				return "", err
			}
			names = append(names, name)
		}
		return names[0], nil
	}
}

const (
	__port = 8000
)

var (
	port         = fmt.Sprintf(":%d", __port)
	database_url = "postgresql://root@cockroach:26257/defaultdb?sslmode=disable"
	cockroachDB  *DataStore // make visible outside of init
	err          error
)

func init() {
	if cockroachDB, err = NewDataStore(database_url); err != nil {
		log.Fatalf("Could not connect to database: %s \n", err.Error())
	} else {
		log.Printf("Database connected \n")
		if err = cockroachDB.CreateTestDatabase(context.Background()); err != nil {
			log.Printf("Error creating database: %s \n", err.Error())
		} else if err = cockroachDB.InsertTestElement(context.Background()); err != nil {
			log.Printf("Error inserting test elements: %s \n", err.Error())
		} else if name, err := cockroachDB.QueryTestElement(context.Background()); err != nil {
			log.Printf("Error querying test elements: %s \n", err.Error())
		} else {
			log.Println(name)
		}
	}
}

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
	api.GET("/data", func(c *gin.Context) {
		res, err := cockroachDB.QueryTestElement(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, res)
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Printf("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown:%s \n", err.Error())
	}

	// catching ctx.Done() after 5 second timeout
	select {
	case <-ctx.Done():
		log.Printf("Server exiting...\n")
	}
	log.Println("Server Shutdown")
}
