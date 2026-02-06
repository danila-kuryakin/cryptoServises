package repository

import (
	"controller/pkg/models"
	"database/sql"
)

// ProposalRepo - работает с Proposals и смежными таблицами
type ProposalRepo interface {
	ReadEvents() ([]models.ProposalEvent, error)
	ProposalDeliverySuccessful(proposals []models.ProposalEvent) error
	EventDeliverySuccessful(event []models.CurrentEvent) error
	AddEventScheduler(proposals []models.ProposalEvent) error
	GetCurrentEvents(number int64) ([]models.CurrentEvent, error)
}

type SpaceRepo interface {
	ReadEvents() ([]models.SpaceEvent, error)
	DeliverySuccessful(proposals []models.SpaceEvent) error
	AddEventScheduler(proposals []models.SpaceEvent) error
	GetCurrentEvents(number int64) ([]models.CurrentEvent, error)
}

type UserRepo interface {
	GetUserSubscriptions() ([]int64, error)
}

type Repository struct {
	ProposalRepo
	SpaceRepo
	UserRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ProposalRepo: NewProposalPostgres(db),
		SpaceRepo:    NewSpacePostgres(db),
		UserRepo:     NewUserPostgres(db),
	}
}
