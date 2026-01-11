package models

type NewData struct {
	TableName string   `json:"table_name"`
	Ids       []string `json:"ids"`
}
