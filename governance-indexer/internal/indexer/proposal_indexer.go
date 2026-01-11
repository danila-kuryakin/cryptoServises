package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"governance-indexer/internal/repository"
	"governance-indexer/pkg/models"
	"governance-indexer/pkg/service"
	"io"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

type ProposalIndexer struct {
	repo    *repository.Repository
	rwKafka *service.ReaderWriterService
}

func NewProposalIndexer(repo *repository.Repository, rwKafka *service.ReaderWriterService) *ProposalIndexer {
	return &ProposalIndexer{repo: repo, rwKafka: rwKafka}
}

var endpoint = "https://hub.snapshot.org/graphql"

// GraphQL-запрос
var queryProposals = `
{
	proposals (
		first: %d,
		skip: 0,
		orderDirection: desc
	) {
		id
		title
		author
		created
		start
		end
		state
		snapshot
		choices
		space {
			id
			name
		}
	}
}`

// CreatedResponse и DataResponse Структуры для получения ответа
type CreatedResponse struct {
	Data DataResponse `json:"data"`
}

// CreatedResponse и DataResponse Структуры для получения ответа
type DataResponse struct {
	Proposals []models.Proposals `json:"proposals"`
}

// IndexProposal получает записи proposal и сохраняет в БД
func (p *ProposalIndexer) IndexProposal(numberRecords int) error {

	// Дописываем запрос. Добавляем количество получаемых записей
	query := fmt.Sprintf(queryProposals, numberRecords)

	// Переводим в json формат
	jsonData, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		log.Println("JSON marshal customErrors:", err)
		return err
	}

	//log.Println("1")
	// Отправляем запрос и получаем ответ
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP request customErrors:", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	//log.Println("2")
	// получаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body customErrors:", err)
		return err
	}

	//log.Println("3")
	// парсим в json
	var result CreatedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("JSON unmarshal customErrors:", err)
		return err
	}

	//log.Println("4")
	// смотрим каких записей нет
	missing, err := p.repo.FindMissing(result.Data.Proposals)
	if err != nil {
		log.Println("Error finding missing proposals:", err)
		return err
	}

	//log.Println("5")
	lenMissing := len(missing)
	if lenMissing == 0 {
		log.Println("No missing proposals found")
		return nil
	}

	// TODO: здесь должен быть outbox-паттерн. Нужно написать часть проверяет доставку сообщений и
	// отправляет запрос заново, если подтверждения нет
	ids := make([]string, 0, lenMissing)
	for _, proposal := range missing {
		ids = append(ids, proposal.ID)
	}

	eventData := models.NewData{
		TableName: "proposals",
		Ids:       ids,
	}

	fmt.Println(eventData)
	data, err := json.Marshal(eventData)
	if err != nil {
		log.Println(fmt.Sprintf("Marshal customErrors: %v", err))
	}
	err = p.rwKafka.WriteMessage(
		kafka.Message{
			Value: data,
		})
	if err != nil {
		log.Println(fmt.Sprintf("Write messages customErrors: %v", err))
	}

	log.Println("Proposals writes:", len(missing))

	// если есть записи, то сохраняем в БД
	if len(missing) > 0 {
		if err := p.repo.AddProposal(missing); err != nil {
			log.Println("Repository customErrors:", err)
			return err
		}
	}

	//log.Println("6")
	return nil
}
