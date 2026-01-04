package service

import (
	"telegramBot/internal/models"
	"telegramBot/internal/repository"
)

type Dao interface {
	GetLastProposals() ([]models.Proposal, error)
}

type Service struct {
	Dao
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Dao: NewDaoService(repo),
	}
}
