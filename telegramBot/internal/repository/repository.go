package repository

import (
	"database/sql"
	"telegramBot/internal/models"
)

type DaoRepo interface {
	GetLastProposals() ([]models.Proposal, error)
}

type Repository struct {
	DaoRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DaoRepo: NewDaoPostgres(db),
	}
}
