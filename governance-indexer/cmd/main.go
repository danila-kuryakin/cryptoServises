package main

import (
	"fmt"
	"governance-indexer/internal/repository"
	"governance-indexer/internal/timer"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"governance-indexer/internal/config"
	"governance-indexer/internal/indexer"
)

func main() {
	// загружаем настройки из .env и yml файлов
	config.LoadEnv(".env")
	cfg := config.LoadConfig("configs/config.yml")

	// Конфигурация и подключение к PostgreSQL
	postgresConf := repository.PostgresConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     cfg.Database.Host,
		Port:     strconv.Itoa(cfg.Database.Port),
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}
	db, err := repository.NewPostgresDB(postgresConf)
	if err != nil {
		log.Fatal(err)
	}

	// Подключения модулей
	repo := repository.NewRepository(db)
	index := indexer.NewIndexer(repo)
	tm := timer.NewTimer(index, cfg)
	go tm.StartProposal()

	log.Println(fmt.Sprintf("Server started on: %s", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil); err != nil {
		return
	}
}
