package service

import (
	"controller/pkg/models"
	pkgService "controller/pkg/service"
	"telegramBot/internal/config"
	"telegramBot/internal/repository"
)

type Dao interface {
	NewUser(userId int64) error
	SubscribedSpaces(userId int64) (bool, error)
	UnsubscribedSpaces(userId int64) (bool, error)
	SubscribedProposals(userId int64) (bool, error)
	UnsubscribedProposals(userId int64) (bool, error)
	StatusSubscribedSpaces(userId int64) (int, error)
	StatusSubscribedProposals(userId int64) (int, error)
	CreateVotesId(userId int64, votesId string) (bool, error)
	DropVotesId(userId int64, votesId string) (bool, error)
	GetVotesByUser(userId int64) ([]string, error)

	KafkaListen() (models.CurrentProposalEvent, error)
}

type Service struct {
	Dao
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	controllerKafka := pkgService.NewReaderWriterService(
		cfg.Kafka.Address,
		cfg.Kafka.Port,
		config.DaoControllerBotTopic,
		config.DaoControllerBotGroup)
	return &Service{
		Dao: NewDaoService(repo, controllerKafka),
	}
}
