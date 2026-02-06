package repository

import (
	customError "controller/errors"
	"controller/pkg/models"
	"controller/pkg/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type SpacePostgres struct {
	db *sql.DB
}

func NewSpacePostgres(db *sql.DB) *SpacePostgres {
	return &SpacePostgres{db: db}
}

func (s *SpacePostgres) ReadEvents() ([]models.SpaceEvent, error) {

	query := fmt.Sprintf(`SELECT spc.space_id, spc.created_at				       
				FROM %s AS spc
				LEFT JOIN %s AS evn ON spc.space_id = evn.space_id
				WHERE evn.processed_at IS NULL`, repository.SpacesTable, repository.SpacesOutboxTable)

	rows, err := s.db.Query(query)
	if err != nil {
		log.Println("Query:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}(rows)

	var modelArr []models.SpaceEvent

	for rows.Next() {
		var p models.SpaceEvent
		if err := rows.Scan(&p.ID, &p.Created); err != nil {
			return nil, err
		}
		modelArr = append(modelArr, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return modelArr, nil
}

func (s *SpacePostgres) DeliverySuccessful(proposals []models.SpaceEvent) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Error in Begin:", err)
		return err
	}
	for _, proposal := range proposals {
		query := fmt.Sprintf(`
				UPDATE %s
				SET processed_at = NOW()
				WHERE hex_id = $1;
				`, repository.ProposalsTable)

		_, err := tx.Exec(query, proposal.ID)
		if err != nil {
			log.Println("Error in Exec:", err)
			return tx.Rollback()
		}
	}

	return tx.Commit()
}

const (
	eventCreateSpace = "create space"
)

func (s *SpacePostgres) AddEventScheduler(spaces []models.SpaceEvent) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Error in Begin:", err)
		return err
	}
	for _, space := range spaces {
		query := fmt.Sprintf(`
				INSERT INTO %s (hex_id, event_type, event_at)
				VALUES ($1, $2, $3)
			`, repository.EventSchedulerTable)

		if space.Created.Valid {
			if _, err := tx.Exec(query, space.ID, eventCreateSpace, space.Created); err != nil {
				log.Println("Error in Exec Created:", err)
				return tx.Rollback()
			}
		}
	}
	return tx.Commit()
}

func (s *SpacePostgres) GetCurrentEvents(number int64) ([]models.CurrentEvent, error) {

	query := fmt.Sprintf(`
			SELECT evn.hex_id, evn.event_type, evn.event_at, spc.id, spc.name, spc.about
			FROM %s AS evn
			LEFT JOIN %s AS spc ON evn.hex_id = spc.space_id
			WHERE evn.processed_at IS NULL
			ORDER BY evn.event_at ASC
			LIMIT $1
		`, repository.EventSchedulerTable, repository.SpacesTable)

	rows, err := s.db.Query(query, number)
	if err != nil {
		log.Println("Queryw:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}(rows)

	var event []models.CurrentEvent

	for rows.Next() {
		var c models.CurrentEvent
		if err := rows.Scan(&c.ID, &c.EventType, &c.EventTime, &c.SpaceID, &c.SpaceName, &c.Title); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, customError.ErrDataNotFound
			}
			return nil, err
		}
		event = append(event, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return event, nil
}
