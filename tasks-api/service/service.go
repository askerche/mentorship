package service

import (
	"context"
	"tasks-api/models"
	"tasks-api/repo"
)

func New(r *repo.Repo) *Service {
	return &Service{
		r: r,
	}
}

type Service struct {
	r *repo.Repo
}

func (s *Service) SaveTask(ctx context.Context, task models.Task) (int, error) {
	id, err := s.r.CreateTask(ctx, task)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) GetTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.r.SelectTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) GetTask(ctx context.Context, taskId int) (models.Task, error) {
	var task models.Task
	task, err := s.r.SelectTask(ctx, taskId)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, taskId int) error {
	err := s.r.DeleteTask(ctx, taskId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateTask(ctx context.Context, taskId int, task models.UpdateTask) error {
	err := s.r.UpdateTask(ctx, task, taskId)
	if err != nil {
		return err
	}
	return nil
}
