package repository

import (
	"database/sql"
)

type PostgresConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// PostgresDSN (Data Source Name) формирует строку подключения PostgreSQL
func PostgresDSN(config PostgresConfig) string {
	return "postgres://" +
		config.UserName + ":" + config.Password + "@" +
		config.Host + ":" + config.Port + "/" +
		config.Name + "?" + "sslmode=" + config.SSLMode
}

func NewPostgresDB(config PostgresConfig) (*sql.DB, error) {
	dsn := PostgresDSN(config)
	//fmt.Println("dsn:", dsn)
	db, err := sql.Open("postgres", dsn)
	//fmt.Println(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
