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
	"time"
)

type ProposalIndexer struct {
	repo *repository.Repository
}

func NewProposalIndexer(repo *repository.Repository) *ProposalIndexer {
	return &ProposalIndexer{repo: repo}
}

var endpoint = "https://hub.snapshot.org/graphql"

// GraphQL-запрос
var createdQuery = `
{
  proposals(first: 5, orderBy: "created", orderDirection: desc) {
    id
    title
    author
    created
    state
    space {
      id
      name
    }
  }
}`

// Структуры для парсинга ответа
type CreatedResponse struct {
	Data DataResponse `json:"data"`
}

type DataResponse struct {
	Proposals []models.Proposals `json:"proposals"`
}

func (proposalIndexer *ProposalIndexer) CreateProposal() error {
	jsonData, err := json.Marshal(map[string]string{
		"query": createdQuery,
	})
	if err != nil {
		log.Println("JSON marshal error:", err)
		return err
	}

	start := time.Now()

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

	var result CreatedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("JSON unmarshal error:", err)
		return err
	}

	elapsed := time.Since(start)

	if err := proposalIndexer.repo.AddProposal(result.Data.Proposals); err != nil {
		return err
	}

	fmt.Printf("\n⏱ Запрос выполнен за: %s\n", elapsed)
	fmt.Println("Последние 5 proposals Snapshot:")
	return nil
}
