package repository

const (
	ProposalsTable       = "proposals"
	ProposalsOutboxTable = "proposals_outbox"
	SpacesTable          = "spaces"
	SpacesOutboxTable    = "spaces_outbox"
	EventSchedulerTable  = "event_scheduler"
	UserTable            = "users"
	UserVotesTable       = "users_votes"
)

const (
	EventProposalCreated = "proposalCreated"
	EventHistory         = "eventHistory"
	EventSpaceCreated    = "spaceCreated"
)
