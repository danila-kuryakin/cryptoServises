package models

//
//// EventType тип события DAO
//type EventType string
//
//const (
//	ProposalCreated  EventType = "PROPOSAL_CREATED"
//	VotingStarted    EventType = "VOTING_STARTED"
//	VoteCast         EventType = "VOTE_CAST"
//	ProposalExecuted EventType = "PROPOSAL_EXECUTED"
//)
//
//// Event структура нормализованного события DAO
//type Event struct {
//	DAO         string      `json:"dao"`         // Название DAO
//	EventType   EventType   `json:"eventType"`   // Тип события
//	Payload     interface{} `json:"payload"`     // Данные события
//	BlockNumber int64       `json:"blockNumber"` // Номер блока (для блокчейн-событий)
//}
