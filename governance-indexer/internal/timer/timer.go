package timer

import (
	"governance-indexer/internal/indexer"
)

type ProposalTimerInterface interface {
	StartProposal()
}

type Timer struct {
	ProposalTimerInterface
}

func NewTimer(index *indexer.Indexer) *Timer {
	return &Timer{
		ProposalTimerInterface: NewProposalTimer(index),
	}
}
