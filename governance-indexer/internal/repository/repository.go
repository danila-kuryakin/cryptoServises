package repository

import (
	"database/sql"
	"governance-indexer/internal/models"
)

type ProposalRepo interface {
	AddProposal(proposals []models.Proposals) error
	//GetDiff(proposals []models.Proposals) ([]models.Proposals, error)
	//GetHexId() ([]ProposalsHexId, error)
}

type Repository struct {
	ProposalRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
	}
}
