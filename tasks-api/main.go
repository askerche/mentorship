package main

import (
	"context"
	"fmt"
	"tasks-api/handlers"
	"tasks-api/repo"
	"tasks-api/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://postgres:mysecretpassword@tasks-postgres:5432/postgres"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		fmt.Println("error connect to db: ", err)
		return
	}
	err = conn.Ping(ctx)
	if err != nil {
		fmt.Println("error ping db: ", err)
		return
	}

	repo := repo.New(conn)
	svc := service.New(repo)
	handler := handlers.New(svc)

	router := gin.Default()
	router.GET("/tasks", handler.GetTasksHandler)
	router.GET("/tasks/:id", handler.GetTaskHandler)
	router.POST("/tasks", handler.CreateTaskHandler)
	router.DELETE("/tasks/:id", handler.DeleteTaskHandler)
	router.PUT("/tasks/:id", handler.UpdateTaskHandler)

	fmt.Println("Запущен")
	router.Run(":8080")
}
