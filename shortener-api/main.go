package main

import (
	"context"
	"fmt"
	"os"

	"shortener/handlers"
	"shortener/repo"
	"shortener/service"

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

	router := gin.Default()

	repo := repo.Repo{
		DB: conn,
	}
	service := service.Service{
		Repo: repo,
	}
	handler := handlers.Handler{
		Svc: service,
	}

	router.Use(cors.Default())
	router.POST("/shorten", handler.ShortenLinkHandler)
	router.GET("/analytics/:link", handler.LinkAnalyticsHandler)
	router.GET("/:link", handler.ShortLinkHandler)

	fmt.Println("Сервер запущен")
	router.Run(":8282")

}
