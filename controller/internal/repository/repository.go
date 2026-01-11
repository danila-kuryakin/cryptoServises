package repository

import (
	"controller/pkg/models"
	"database/sql"
)

// ProposalRepo - работает с Proposals и смежными таблицами
type ProposalRepo interface {
	ReadProposalsEvents() ([]models.ProposalEvent, error)
	ProposalDeliverySuccessful(proposals []models.ProposalEvent) error
	EventDeliverySuccessful(event []models.CurrentEvent) error
	AddEventScheduler(proposals []models.ProposalEvent) error
	GetCurrentEvents(number int64) ([]models.CurrentEvent, error)
}

type UserRepo interface {
	GetUserSubscriptions() ([]int64, error)
}

type Repository struct {
	ProposalRepo
	UserRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
		UserRepo:     NewUserPostgres(db),
	}
}
