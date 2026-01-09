package service

import (
	"controller/internal/repository"
	"controller/pkg/models"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type DaoService struct {
	repo *repository.Repository
}

func NewDaoService(repo *repository.Repository) *DaoService {
	return &DaoService{repo: repo}
}

func (d *DaoService) Processing(message *kafka.Message) error {
	var proposals_in []models.Proposals
	if err := json.Unmarshal(message.Value, &proposals_in); err != nil {
		log.Println("bad json:", err)
	}

	//proposals_in, err := d.repo.ReadNewProposals()
	//if err != nil {
	//	log.Println("Error reading proposals: ", err)
	//	return err
	//}

	if err := d.repo.DeliverySuccessful(proposals_in); err != nil {
		log.Println("delivery error:", err)
		return err
	}

	for i, proposal := range proposals_in {
		log.Printf("  %d) received order: %#v\n", i, proposal)
	}
	return nil
}
