package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"telegramBot/internal/bot"
	"telegramBot/internal/config"
	"telegramBot/internal/repository"
	"telegramBot/internal/service"
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
		log.Println(err)
	}

	// Подключения модулей
	repo := repository.NewRepository(db)
	services := service.NewService(repo, cfg)
	bots := bot.NewBot(services, cfg)
	go bots.StartBot()

	// Создаем http сервер
	log.Println(fmt.Sprintf("Server started on: %s", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil); err != nil {
		return
	}
}
