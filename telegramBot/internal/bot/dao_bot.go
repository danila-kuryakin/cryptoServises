package bot

import (
	"fmt"
	"log"
	"os"
	"telegramBot/internal/config"
	"telegramBot/internal/service"
	"time"

	"gopkg.in/telebot.v3"
)

type DaoBot struct {
	service *service.Service
	config  *config.Config
}

func NewDaoBot(service *service.Service, config *config.Config) *DaoBot {
	return &DaoBot{service: service, config: config}
}

func (p DaoBot) StartBot() {

	bot_token := os.Getenv("API_KEY")

	pref := telebot.Settings{
		Token:  bot_token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// Команда /start
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Я родился! Напиши что нибудь.")
	})

	// Эхо на любые текстовые сообщения
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		proposals, err := p.service.GetLastProposals()
		if err != nil {
			return err
		}
		retString := "Вот последние 5 proposals:\n"
		for _, p := range proposals {

			retString += fmt.Sprintf(
				" Created: %s\n Start: %s\n End: %s\n Title: %s\n Space: %s (%s)\n Author: %s\n State: %s\n\n",
				p.Created.Format("2006-01-02 15:04:05"),
				p.Start.Format("2006-01-02 15:04:05"),
				p.End.Format("2006-01-02 15:04:05"),
				p.Title,
				p.SpaceName,
				p.SpaceId,
				p.Author,
				p.State,
			)
		}
		return c.Send(retString)
	})

	bot.Start()

}
