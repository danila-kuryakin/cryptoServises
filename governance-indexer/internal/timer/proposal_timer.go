package timer

import (
	"governance-indexer/internal/config"
	"governance-indexer/internal/indexer"
	"time"
)

type ProposalTimer struct {
	index  *indexer.Indexer
	config *config.Config
}

func NewProposalTimer(index *indexer.Indexer, config *config.Config) *ProposalTimer {
	return &ProposalTimer{index: index, config: config}
}

func (p ProposalTimer) StartProposal() {
	for {
		err := p.index.IndexProposal(p.config.Proposal.NumberRecords)
		if err != nil {
			return
		}

		var durationMinutes = time.Duration(p.config.Proposal.TimeRequest * 60)

		time.Sleep(durationMinutes * time.Second)
	}
}
