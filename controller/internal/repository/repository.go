package repository

import (
	"controller/pkg/models"
	"database/sql"
)

type ProposalRepo interface {
	ReadNewProposals() ([]models.Proposals, error)
	DeliverySuccessful(proposals []models.Proposals) error
}

type Repository struct {
	ProposalRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
	}
}
