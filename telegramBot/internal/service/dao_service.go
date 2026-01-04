package service

import (
	"telegramBot/internal/models"
	"telegramBot/internal/repository"
)

type DaoService struct {
	repo *repository.Repository
}

func NewDaoService(repo *repository.Repository) *DaoService {
	return &DaoService{repo: repo}
}

func (p DaoService) GetLastProposals() ([]models.Proposal, error) {
	return p.repo.GetLastProposals()
}
