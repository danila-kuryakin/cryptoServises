package models

import (
	"database/sql"
)

type NewData struct {
	TableName string   `json:"table_name"`
	IDs       []string `json:"ids"`
}

type ProposalEvent struct {
	ID      string       `db:"hex_id"`  // Уникальный идентификатор пропозиции
	Created sql.NullTime `db:"created"` // Время создания записи
	Start   sql.NullTime `db:"start"`   // Время начала голосования
	End     sql.NullTime `db:"end"`     // Время окончание голосования
}

type CurrentEvent struct {
	ID        string       `db:"hex_id"`
	EventType string       `db:"event_type"`
	EventTime sql.NullTime `db:"event_at"`
	SpaceID   string       `db:"space_id"`
	SpaceName string       `db:"space_name"`
	Title     string       `db:"title"`
}

type CurrentProposalEvent struct {
	Users        []int64      `json:"user_id"`
	CurrentEvent CurrentEvent `json:"current_event"`
}
