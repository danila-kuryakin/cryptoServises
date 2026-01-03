package repository

import (
	"database/sql"
)

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// PostgresDSN (Data Source Name) формирует строку подключения PostgreSQL
func PostgresDSN(config PostgresConfig) string {
	return "postgres://" +
		config.Username + ":" + config.Password + "@" +
		config.Host + ":" + config.Port + "/" +
		config.Name + "?" + "sslmode=" + config.SSLMode
}

func NewPostgresDB(config PostgresConfig) (*sql.DB, error) {
	dsn := PostgresDSN(config)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
