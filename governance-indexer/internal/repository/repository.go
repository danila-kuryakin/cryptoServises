package repository

import (
	"database/sql"
	"governance-indexer/internal/models"
)

type ProposalRepo interface {
	AddProposal(proposals []models.Proposals) error
	FindMissing(proposals []models.Proposals) ([]models.Proposals, error)
}

type Repository struct {
	ProposalRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
	}
}
