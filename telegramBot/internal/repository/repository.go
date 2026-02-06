package repository

import (
	"controller/pkg/models"
	"database/sql"
)

type UserRepo interface {
	GetUserById(userId int64) (*models.User, error)
	CreateUser(userId int64) error
	SetSubscribedSpaces(userId int64, subscribeStatus int) (bool, error)
	SetSubscribedProposals(userId int64, subscribeStatus int) (bool, error)
	StatusSubscribedSpaces(userId int64) (int, error)
	StatusSubscribedProposals(userId int64) (int, error)
	CreateVotesId(userId int64, votesId string) (bool, error)
	DropVotesId(userId int64, votesId string) (bool, error)
	GetVotesByUser(userId int64) ([]string, error)
}

type Repository struct {
	UserRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepo: NewUserPostgres(db),
	}
}
