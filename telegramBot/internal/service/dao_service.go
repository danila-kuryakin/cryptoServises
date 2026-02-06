package service

import (
	"controller/pkg/models"
	pkgService "controller/pkg/service"
	"encoding/json"
	"errors"
	"log"
	cusmomError "telegramBot/internal/error"
	"telegramBot/internal/repository"
)

type DaoService struct {
	repo            *repository.Repository
	controllerKafka *pkgService.ReaderWriterService
}

func NewDaoService(repo *repository.Repository, controllerKafka *pkgService.ReaderWriterService) *DaoService {
	return &DaoService{repo: repo, controllerKafka: controllerKafka}
}

// NewUser - создание пользователя по id
func (s *DaoService) NewUser(userId int64) error {
	_, err := s.repo.GetUserById(userId)
	if err != nil {
		if errors.Is(err, cusmomError.ErrorUserNotFound) {
			err = s.repo.CreateUser(userId)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

// SubscribedSpaces - подписать пользователя на Spaces
func (s *DaoService) SubscribedSpaces(userId int64) (bool, error) {
	_, err := s.repo.SetSubscribedSpaces(userId, 1)
	if err != nil {
		return false, err
	}
	log.Println("Subscribed", userId)
	return true, nil
}

// UnsubscribedSpaces - отписать пользователя на Spaces
func (s *DaoService) UnsubscribedSpaces(userId int64) (bool, error) {
	_, err := s.repo.SetSubscribedSpaces(userId, 0)
	if err != nil {
		return false, err
	}
	log.Println("Unsubscribed", userId)
	return true, nil
}

// SubscribedProposals - подписать пользователя на Proposals
func (s *DaoService) SubscribedProposals(userId int64) (bool, error) {
	_, err := s.repo.SetSubscribedProposals(userId, 1)
	if err != nil {
		return false, err
	}
	log.Println("Subscribed", userId)
	return true, nil
}

// UnsubscribedProposals - отписать пользователя на Proposals
func (s *DaoService) UnsubscribedProposals(userId int64) (bool, error) {
	_, err := s.repo.SetSubscribedProposals(userId, 0)
	if err != nil {
		return false, err
	}
	log.Println("Unsubscribed", userId)
	return true, nil
}

func (s *DaoService) StatusSubscribedSpaces(userId int64) (int, error) {
	return s.repo.StatusSubscribedSpaces(userId)
}

func (s *DaoService) StatusSubscribedProposals(userId int64) (int, error) {
	return s.repo.StatusSubscribedProposals(userId)
}

func (s *DaoService) KafkaListen() (models.CurrentProposalEvent, error) {

	message, err := s.controllerKafka.ReadMessage()
	if err != nil {
		log.Println("Error reading message", err)
		return models.CurrentProposalEvent{}, err
	}

	var eventData models.CurrentProposalEvent
	if err := json.Unmarshal(message.Value, &eventData); err != nil {
		log.Println("bad json:", err)
	}

	return eventData, nil
}

func (s *DaoService) CreateVotesId(userId int64, votesId string) (bool, error) {
	return s.repo.CreateVotesId(userId, votesId)
}

func (s *DaoService) DropVotesId(userId int64, votesId string) (bool, error) {
	return s.repo.DropVotesId(userId, votesId)
}

func (s *DaoService) GetVotesByUser(userId int64) ([]string, error) {
	return s.repo.GetVotesByUser(userId)
}
