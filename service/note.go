package service

import (
	"context"

	"github.com/viveknathani/nattukaka/entity"
)

func (service *Service) CreateNote(ctx context.Context, n *entity.Note) error {

	service.Logger.Info("database: insert note start.", zapReqID(ctx))
	err := service.Repo.CreateNote(n)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoInsert
	}
	service.Logger.Info("database: insert note end.", zapReqID(ctx))
	return nil
}

func (service *Service) UpdateNote(ctx context.Context, n *entity.Note) error {

	service.Logger.Info("database: update note start.", zapReqID(ctx))
	err := service.Repo.UpdateNote(n)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoInsert
	}
	service.Logger.Info("database: update note end.", zapReqID(ctx))
	return nil
}

func (service *Service) GetAllNotes(ctx context.Context, userId string) (*[]entity.Note, error) {

	service.Logger.Info("database: fetch notes start.", zapReqID(ctx))
	notes, err := service.Repo.GetAllNotes(userId)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return nil, ErrNoFetch
	}
	service.Logger.Info("database: fetch notes end.", zapReqID(ctx))
	return notes, nil
}

func (service *Service) GetNote(ctx context.Context, id string, userId string) (*[]entity.Note, error) {

	service.Logger.Info("database: fetch note start.", zapReqID(ctx))
	note, err := service.Repo.GetNote(id, userId)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return nil, ErrNoFetch
	}
	service.Logger.Info("database: fetch note end.", zapReqID(ctx))
	return note, nil
}
