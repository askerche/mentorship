package main

import (
	"context"
	"fmt"
	"tasks-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	connString := "postgres://postgres:mysecretpassword@localhost:5433/postgres"
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

	Handler := handlers.Handler{
		DB: conn,
	}

	router := gin.Default()
	router.GET("/tasks", Handler.GetTasksHandler)
	router.GET("/tasks/:id", Handler.GetTaskHandler)
	router.POST("/tasks", Handler.CreateTaskHandler)
	router.DELETE("/tasks/:id", Handler.DeleteTaskHandler)
	router.PUT("/tasks/:id", Handler.UpdateTaskHandler)

	fmt.Println("Запущен")
	router.Run(":8181")
}
