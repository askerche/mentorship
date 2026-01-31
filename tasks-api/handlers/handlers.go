package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"tasks-api/models"
	"tasks-api/service.go"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	svc service.Service
}

func (h *Handler) CreateTaskHandler(c *gin.Context) {
	var task models.Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		fmt.Println("error unmarhall: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	id, err := h.svc.SaveTask(c, task)
	if err != nil {
		fmt.Println("error insert to db: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed saving tasks"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Status": "Task created!",
		"id":     id,
	})
}

func (h *Handler) GetTasksHandler(c *gin.Context) {
	tasks := []models.Task{}
	tasks, err := h.svc.GetTasks(c)
	if err != nil {
		fmt.Println("error query from db: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetTaskHandler(c *gin.Context) {
	taskId := c.Param("id")
	task, err := h.svc.GetTask(c, taskId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "failed to find tasks"})
		} else {
			fmt.Println("db error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		}
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTaskHandler(c *gin.Context) {
	taskId := c.Param("id")
	res, err := h.svc.DeleteTask(c, taskId)
	if err != nil {
		fmt.Println("error delete task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Task deleted."})
}

func (h *Handler) UpdateTaskHandler(c *gin.Context) {
	taskId := c.Param("id")
	var task models.UpdateTask
	err := c.ShouldBindJSON(&task)
	if err != nil {
		fmt.Println("error unmarhall update task: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err = h.svc.UpdateTask(c, taskId, task)
	if err != nil {
		fmt.Println("error update task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, "Task successfully update.")
}
