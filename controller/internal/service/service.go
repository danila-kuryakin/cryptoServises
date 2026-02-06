package service

import (
	"controller/internal/config"
	"controller/internal/repository"
	"controller/pkg/service"

	"github.com/segmentio/kafka-go"
)

type Dao interface {
	Processing() error
	CompareIDs(a, b []string) (same, onlyA, onlyB []string)
	ProcessingProposals(ids []string) error
	ProcessingSpaces(ids []string) error
	MessageControllerProposal() error
	MessageControllerSpace() error
}

type ReaderWriterKafka interface {
	ReadMessage() (*kafka.Message, error)
	WriteMessage(message kafka.Message) error
}

type Service struct {
	Dao
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	indexerKafka := service.NewReaderWriterService(cfg.Kafka.Address, cfg.Kafka.Port, config.DaoIndexerTopic, config.DaoIndexerGroup)
	botKafka := service.NewReaderWriterService(cfg.Kafka.Address, cfg.Kafka.Port, config.DaoControllerBotTopic, config.DaoControllerBotGroup)

	return &Service{
		Dao: NewDaoService(repo, indexerKafka, botKafka),
	}
}
