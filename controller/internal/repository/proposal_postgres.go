package repository

import (
	"controller/pkg/models"
	"database/sql"
	"fmt"
	"log"
)

type ProposalPostgres struct {
	db *sql.DB
}

func NewProposalPostgres(db *sql.DB) *ProposalPostgres {
	return &ProposalPostgres{db: db}
}

func (p *ProposalPostgres) ReadNewProposals() ([]models.Proposals, error) {

	query := fmt.Sprintf(`
				SELECT prop.hex_id AS id, prop.title, prop.author, prop.created_at, prop.start_at, prop.end_at, 
				       prop.snapsho, prop.state, prop.choices, prop.space_id, prop.space_name
				FROM %s AS prop
				LEFT JOIN %s AS evn ON prop.id = evn.hex_id
				WHERE evn.processed_at IS NULL;
				`, proposalsTable, eventOutboxTable)

	rows, err := p.db.Query(query)
	if err != nil {
		log.Println("Query:", err)
		return nil, err
	}
	defer rows.Close()

	var proposals []models.Proposals

	for rows.Next() {
		var p models.Proposals
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		proposals = append(proposals, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return proposals, nil
	return nil, nil
}

func (p *ProposalPostgres) DeliverySuccessful(proposals []models.Proposals) error {
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Error in Begin:", err)
		return err
	}
	for _, proposal := range proposals {
		query := fmt.Sprintf(`
				UPDATE %s
				SET processed_at = NOW()
				WHERE hex_id = $1;
				`, eventOutboxTable)

		_, err := tx.Exec(query, proposal.ID)
		if err != nil {
			log.Println("Error in Exec:", err)
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
