package service

import (
	"context"
	"fmt"
	"tasks-api/models"
	"tasks-api/repo"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	r repo.Repo
}

func (s *Service) SaveTask(ctx context.Context, task models.Task) (int, error) {
	var id int
	err := s.r.DB.QueryRow(ctx, `INSERT INTO tasks (title, description, status)
	 VALUES($1, $2, COALESCE(NULLIF($3, ''), 'new')) RETURNING id`, task.Title, task.Description, task.Status).Scan(&id)
	if err != nil {
		fmt.Println("error insert to db: ", err)
		return 0, err
	}
	return id, nil
}

func (s *Service) GetTasks(ctx context.Context) ([]models.Task, error) {
	tasks := []models.Task{}
	rows, err := s.r.DB.Query(ctx, `SELECT id, title, description, status, created_at FROM tasks`)
	if err != nil {
		fmt.Println("error query from db: ", err)
		return []models.Task{}, err
	}
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			fmt.Println("error scan from db: ", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *Service) GetTask(ctx context.Context, taskId string) (models.Task, error) {
	var task models.Task
	err := s.r.DB.QueryRow(ctx, `SELECT id, title, description, status, created_at FROM tasks WHERE id = $1`, taskId).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, taskId string) (pgconn.CommandTag, error) {
	res, err := s.r.DB.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, taskId)
	if err != nil {
		fmt.Println("error delete task: ", err)
		return pgconn.CommandTag{}, err
	}
	return res, nil
}

func (s *Service) UpdateTask(ctx context.Context, taskId string, task models.UpdateTask) error {
	res, err := s.r.DB.Exec(ctx, `UPDATE tasks SET title = coalesce($1, title),
	 description = coalesce($2, description), status = coalesce($3, status) WHERE ID = $4`,
		task.Title, task.Description, task.Status, taskId,
	)
	if err != nil {
		fmt.Println("error update task: ", err)
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("failed to find task id")
	}
	return nil
}
