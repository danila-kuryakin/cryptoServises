package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"governance-indexer/pkg/models"
	"log"
	"time"

	"github.com/lib/pq"
)

type ProposalPostgres struct {
	db *sql.DB
}

func NewProposalPostgres(db *sql.DB) *ProposalPostgres {
	return &ProposalPostgres{db: db}
}

// AddProposal Добавляет новые proposals в БД
func (p ProposalPostgres) AddProposal(proposals []models.Proposals) error {
	if len(proposals) == 0 {
		log.Println("AddProposal called with no proposals")
		return nil
	}

	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Error in proposalPostgres.AddProposal:", err)
		return err
	}

	for _, proposal := range proposals {
		choicesJSON, err := json.Marshal(proposal.Choices)
		if err != nil {
			log.Println("AddProposal err:", err)
			return err
		}

		query := fmt.Sprintf(`
				INSERT INTO %s (hex_id, title, author, created_at, start_at, end_at, 
				                snapshot, state, choices, space_id, space_name)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (hex_id) DO NOTHING
			`, proposalsTable)

		queryOutbox := fmt.Sprintf(`
				INSERT INTO %s (hex_id, event_type, created_at)
				VALUES ($1, $2, now()) ON CONFLICT (hex_id) DO NOTHING
			`, eventOutboxTable)

		_, err = tx.Exec(query,
			proposal.ID,
			proposal.Title,
			proposal.Author,
			time.Unix(proposal.Created, 0),
			time.Unix(proposal.Start, 0),
			time.Unix(proposal.End, 0),
			proposal.Snapshot,
			proposal.State,
			string(choicesJSON),
			proposal.Space.ID,
			proposal.Space.Name)
		if err != nil {
			log.Println("Error to exec proposal :", err)
			err := tx.Rollback()
			if err != nil {
				log.Println("Error to Rollback proposal :", err)
				return err
			}
			return err
		}

		_, err = tx.Exec(queryOutbox,
			proposal.ID,
			eventProposalCreated,
		)
		if err != nil {
			log.Println("Error to exec outbox :", err)
			err := tx.Rollback()
			if err != nil {
				log.Println("Error to Rollback outbox :", err)
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

// FindMissing Возвращает proposals которых нет в БД.
// Читает из БД записи из колонки hex_id, количество которых равно proposals в аргументе,
// сравнивает по Proposals.ID (hex_id) каждый элемент
func (p ProposalPostgres) FindMissing(proposals []models.Proposals) ([]models.Proposals, error) {
	// собираем ID из API
	ids := make([]string, 0, len(proposals))
	for _, p := range proposals {
		ids = append(ids, p.ID)
	}
	query := fmt.Sprintf(`
        SELECT unnest($1::text[]) AS hex_id
        EXCEPT
        SELECT hex_id FROM %s;
    `, proposalsTable)
	// база говорит, каких ID у неё нет
	rows, err := p.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	missingIDs := make(map[string]struct{})
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)
		missingIDs[id] = struct{}{}
	}

	// собираем отсутствующие proposals
	var missing []models.Proposals
	for _, p := range proposals {
		if _, ok := missingIDs[p.ID]; ok {
			missing = append(missing, p)
		}
	}

	return missing, nil
}

//// AddProposal Добавляет новые proposals в БД
//func (p ProposalPostgres) AddProposal(proposals []models.Proposals) error {
//	if len(proposals) == 0 {
//		log.Println("AddProposal called with no proposals")
//		return nil
//	}
//
//	log.Println("11")
//	placeholders := make([]string, 0, len(proposals))
//	args := make([]interface{}, 0, len(proposals)*9)
//
//	placeholdersOutbox := make([]string, 0, len(proposals))
//	argsOutbox := make([]interface{}, 0, len(proposals)*3)
//
//	iProposals := 1
//	iOutbox := 1
//	for _, t := range proposals {
//		placeholders = append(placeholders, fmt.Sprintf(""+
//			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
//			iProposals, iProposals+1, iProposals+2, iProposals+3, iProposals+4,
//			iProposals+5, iProposals+6, iProposals+7, iProposals+8,
//			iProposals+9, iProposals+10))
//
//		// choices нужно сериализовать в JSON
//		choicesJSON, err := json.Marshal(t.Choices)
//		if err != nil {
//			log.Println("AddProposal err:", err)
//			return err
//		}
//
//		args = append(
//			args,
//			t.ID,
//			t.Title,
//			t.Author,
//			time.Unix(t.Created, 0),
//			time.Unix(t.Start, 0),
//			time.Unix(t.End, 0),
//			t.Snapshot,
//			t.State,
//			string(choicesJSON),
//			t.Space.ID,
//			t.Space.Name)
//		iProposals += 11
//
//		placeholdersOutbox = append(placeholdersOutbox, fmt.Sprintf(""+
//			"($%d, $%d, now())",
//			iOutbox, iOutbox+1))
//
//		argsOutbox = append(
//			argsOutbox,
//			t.ID,
//			"ProposalCreated",
//		)
//		iOutbox += 2
//	}
//	log.Println("22")
//	query := fmt.Sprintf(`
//      INSERT INTO %s (hex_id, title, author, created_at, start_at, end_at, snapshot, state, choices, space_id, space_name)
//      VALUES %s ON CONFLICT (hex_id) DO NOTHING
//  `, proposalsTable, strings.Join(placeholders, ", "))
//
//	queryOutbox := fmt.Sprintf(`
//      INSERT INTO %s (hex_id, event_type, created_at)
//      VALUES %s ON CONFLICT (hex_id) DO NOTHING
//  `, eventOutboxTable, strings.Join(placeholdersOutbox, ", "))
//	log.Println("33")
//	tx, err := p.db.Begin()
//	if err != nil {
//		log.Println("Error in proposalPostgres.AddProposal:", err)
//		return err
//	}
//
//	fmt.Println(query)
//	fmt.Println(args)
//	fmt.Println(queryOutbox)
//	fmt.Println(argsOutbox)
//
//	log.Println("44")
//	_, err = tx.Exec(query, args...)
//	if err != nil {
//		log.Println("Error to exec proposal :", err)
//		err := tx.Rollback()
//		if err != nil {
//			log.Println("Error to Rollback proposal :", err)
//			return err
//		}
//		return err
//	}
//
//	log.Println("55")
//	_, err = tx.Exec(queryOutbox, argsOutbox...)
//	if err != nil {
//		log.Println("Error to exec outbox :", err)
//		err := tx.Rollback()
//		if err != nil {
//			log.Println("Error to Rollback outbox :", err)
//			return err
//		}
//		return err
//	}
//	log.Println("66")
//	return nil
//}
