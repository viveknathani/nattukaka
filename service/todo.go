package service

import (
	"context"

	"github.com/viveknathani/nattukaka/entity"
)

func (service *Service) CreateTodo(ctx context.Context, t *entity.Todo) error {

	service.Logger.Info("database: insert todo start.", zapReqID(ctx))
	err := service.Repo.CreateTodo(t)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoInsert
	}
	service.Logger.Info("database: insert todo end.", zapReqID(ctx))
	return nil
}

func (service *Service) UpdateTodo(ctx context.Context, t *entity.Todo) error {

	service.Logger.Info("database: update todo start.", zapReqID(ctx))
	err := service.Repo.UpdateTodo(t)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoInsert
	}
	service.Logger.Info("database: update todo end.", zapReqID(ctx))
	return nil
}

func (service *Service) GetAllPendingTodos(ctx context.Context, userId string) (*[]entity.Todo, error) {

	service.Logger.Info("database: fetch todos start.", zapReqID(ctx))
	todos, err := service.Repo.GetAllPendingTodos(userId)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return nil, ErrNoInsert
	}
	service.Logger.Info("database: fetch todos end.", zapReqID(ctx))
	return todos, nil
}

func (service *Service) DeleteTodo(ctx context.Context, id string) error {

	service.Logger.Info("database: delete todo start.", zapReqID(ctx))
	err := service.Repo.DeleteTodo(id)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoRemove
	}
	service.Logger.Info("database: delete todo end.", zapReqID(ctx))
	return nil
}
