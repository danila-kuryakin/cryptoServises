package repository

import (
	"database/sql"
	"encoding/json"
	"telegramBot/internal/models"
)

type DaoPostgres struct {
	db *sql.DB
}

func NewDaoPostgres(db *sql.DB) *DaoPostgres {
	return &DaoPostgres{db: db}
}

func (p *DaoPostgres) GetLastProposals() ([]models.Proposal, error) {
	query := `SELECT 
    		id,
			hex_id,
			title,   
			author, 
			created_at,
			start_at,
			end_at, 
			snapshot,
			state, 
			choices,
			space_id,
			space_name
        FROM proposal
        ORDER BY start_at DESC
        LIMIT 5;`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Proposal

	for rows.Next() {
		var p models.Proposal
		var choicesJSON []byte

		err := rows.Scan(
			&p.ID,
			&p.HexId,
			&p.Title,
			&p.Author,
			&p.Created,
			&p.Start,
			&p.End,
			&p.Snapshot,
			&p.State,
			&choicesJSON,
			&p.SpaceId,
			&p.SpaceName,
		)
		if err != nil {
			return nil, err
		}

		// распарсим JSON -> []string
		if len(choicesJSON) > 0 {
			if err := json.Unmarshal(choicesJSON, &p.Choices); err != nil {
				return nil, err
			}
		}

		result = append(result, p)
	}

	return result, rows.Err()
}
