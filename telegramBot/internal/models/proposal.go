package models

import (
	"time"
)

type Proposal struct {
	ID        int64      `json:"id"`         // Уникальный идентификатор в БД
	HexId     string     `json:"hex_id"`     // Уникальный идентификатор пропозиции
	Title     string     `json:"title"`      // Текст заголовка
	Author    string     `json:"author"`     // Автор (адрес валидного кошелька)
	Created   *time.Time `json:"created_at"` // Время создания записи
	Start     *time.Time `json:"start_at"`   // Время начала голосования
	End       *time.Time `json:"end_at"`     // Время окончание голосования
	Snapshot  int64      `json:"snapshot"`   // Номер блока, на котором базируется голосование
	State     string     `json:"state"`      // Статус (active, closed, pending)
	Choices   []string   `json:"choices"`    // Варианты для голосования
	SpaceId   string     `json:"space_id"`   // Информация по токену
	SpaceName string     `json:"space_name"` // Информация по токену
}
