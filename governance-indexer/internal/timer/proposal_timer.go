package timer

import (
	"governance-indexer/internal/indexer"
	"time"
)

type ProposalTimer struct {
	index *indexer.Indexer
}

func NewProposalTimer(index *indexer.Indexer) *ProposalTimer {
	return &ProposalTimer{index: index}
}

func (p ProposalTimer) StartProposal() {
	for {
		err := p.index.CreateProposal()
		if err != nil {
			return
		}

		time.Sleep(3600 * time.Second)
	}
}
