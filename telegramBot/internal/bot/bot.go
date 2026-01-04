package bot

import (
	"telegramBot/internal/config"
	"telegramBot/internal/service"
)

type DaoBotInterface interface {
	StartBot()
}

type Bot struct {
	DaoBotInterface
}

func NewBot(service *service.Service, cfg *config.Config) *Bot {
	return &Bot{
		DaoBotInterface: NewDaoBot(service, cfg),
	}
}
