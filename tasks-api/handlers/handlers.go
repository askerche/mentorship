package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"tasks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	DB *pgx.Conn
}

func (h *Handler) CreateTaskHandler(c *gin.Context) {
	var task models.Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		fmt.Println("error unmarhall: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var id int
	err = h.DB.QueryRow(c, `INSERT INTO tasks (title, description, status)
	 VALUES($1, $2, COALESCE(NULLIF($3, ''), 'new')) RETURNING id`, task.Title, task.Description, task.Status).Scan(&id)
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
	rows, err := h.DB.Query(c, `SELECT id, title, description, status, created_at FROM tasks`)
	if err != nil {
		fmt.Println("error query from db: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find tasks"})
		return
	}
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			fmt.Println("error scan from db: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find tasks"})
			return
		}
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetTaskHandler(c *gin.Context) {
	taskId := c.Param("id")
	var task models.Task
	err := h.DB.QueryRow(c, `SELECT id, title, description, status, created_at FROM tasks WHERE id = $1`, taskId).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "failed to find tasks"})
			return
		} else {
			fmt.Println("db error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
			return
		}
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTaskHandler(c *gin.Context) {
	taskId := c.Param("id")
	res, err := h.DB.Exec(c, `DELETE FROM tasks WHERE id = $1`, taskId)
	if err != nil {
		fmt.Println("error delete task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, "Task deleted.")
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

	res, err := h.DB.Exec(c, `UPDATE tasks SET title = coalesce($1, title),
	 description = coalesce($2, description), status = coalesce($3, status) WHERE ID = $4`,
		task.Title, task.Description, task.Status, taskId,
	)
	if err != nil {
		fmt.Println("error update task: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	if res.RowsAffected() == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find task id"})
		return
	}
	c.JSON(http.StatusOK, "Task successfully update.")
}
