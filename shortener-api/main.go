package main

import (
	"context"
	"fmt"
	"os"
	"shortener/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	dbUrl := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	h := handlers.New(conn)

	router := gin.Default()

	router.Use(cors.Default())
	router.POST("/shorten", h.ShortenLinkHandler)
	router.GET("/analytics/:link", h.LinkAnalyticsHandler)
	router.GET("/:link", h.ShortLinkHandler)

	fmt.Println("Сервер запущен")
	router.Run(":8282")

}
