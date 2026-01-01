package indexer

import (
	"governance-indexer/internal/repository"
)

type ProposalIndexerInterface interface {
	CreateProposal() error
}

type Indexer struct {
	ProposalIndexerInterface
}

func NewIndexer(repo *repository.Repository) *Indexer {
	return &Indexer{
		ProposalIndexerInterface: NewProposalIndexer(repo),
	}
}
