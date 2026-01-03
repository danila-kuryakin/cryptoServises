package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"governance-indexer/internal/models"
	"governance-indexer/internal/repository"
	"io"
	"log"
	"net/http"
)

type ProposalIndexer struct {
	repo *repository.Repository
}

func NewProposalIndexer(repo *repository.Repository) *ProposalIndexer {
	return &ProposalIndexer{repo: repo}
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

// CreatedResponse и DataResponse Структуры для парсинга ответа
type CreatedResponse struct {
	Data DataResponse `json:"data"`
}

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
		log.Println("JSON marshal error:", err)
		return err
	}
	// Отправляем запрос и получаем ответ
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP request error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body error:", err)
		return err
	}

	// парсим в json
	var result CreatedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("JSON unmarshal error:", err)
		return err
	}

	// смотрим каких записей нет
	missing, err := p.repo.FindMissing(result.Data.Proposals)
	if err != nil {
		fmt.Println("Error finding missing proposals:", err)
		return err
	}
	//for _, proposal := range missing {
	//	fmt.Println(proposal)
	//}
	// TODO: здесь должен быть outbox-паттерн

	// если есть записи, то сохраняем в БД
	if len(missing) > 0 {
		if err := p.repo.AddProposal(result.Data.Proposals); err != nil {
			log.Println("Repository error:", err)
			return err
		}
	}

	return nil
}
