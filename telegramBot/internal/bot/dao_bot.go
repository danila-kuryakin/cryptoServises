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

const (
	subscribedText   = "Подписаться"
	unsubscribedText = "Отменить подписку"
)

const (
	subscribedEvent   = "subscribe_events"
	unsubscribedEvent = "unsubscribe_events"
)

type DaoBot struct {
	service *service.Service
	config  *config.Config
}

func NewDaoBot(service *service.Service, config *config.Config) *DaoBot {
	return &DaoBot{service: service, config: config}
}

func (dao DaoBot) InitListener(bot *telebot.Bot) {
	for {

		listen, err := dao.service.KafkaListen()
		if err != nil {
			return
		}
		text := fmt.Sprintf("id: %s\nEvent type: %s\nDate: %s\nTime: %s\n",
			listen.CurrentEvent.SpaceID,
			listen.CurrentEvent.EventType,
			listen.CurrentEvent.EventTime.Time.Format("2006-01-02"),
			listen.CurrentEvent.EventTime.Time.Format("15:04"),
		)

		for _, user := range listen.Users {
			user := &telebot.User{ID: user}
			_, err := bot.Send(user, text)
			if err != nil {
				return
			}
		}
	}
}

func (dao DaoBot) StartBot() {

	bot_token := os.Getenv("API_KEY")

	pref := telebot.Settings{
		Token:  bot_token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Println(err)
	}

	go dao.InitListener(bot)

	bot.Handle("/start", func(c telebot.Context) error {
		log.Printf("Start| chat: %d, sender: %d, username: %s", c.Chat().ID, c.Sender().ID, c.Sender().Username)
		if err := dao.service.NewUser(c.Chat().ID, c.Sender().Username); err != nil {
			log.Println("Start: ", err)
		}
		// Создаём клавиатуру
		markup := &telebot.ReplyMarkup{ResizeKeyboard: true}
		btnSubscribe := markup.Data(subscribedText, subscribedEvent)

		markup.Inline(
			markup.Row(btnSubscribe),
		)
		return c.Send(
			"Привет! Хочешь получать уведомления о событиях?",
			markup,
		)
	})

	// Обработка нажатия на кнопку
	bot.Handle(&telebot.Btn{Unique: subscribedEvent}, func(c telebot.Context) error {
		log.Printf("Subscribed| chat: %d, sender: %d, username: %s", c.Chat().ID, c.Sender().ID, c.Sender().Username)
		status, err := dao.service.Subscribed(c.Chat().ID)
		if err != nil {
			log.Println("Subscribed: ", err)
			return err
		}

		markup := &telebot.ReplyMarkup{ResizeKeyboard: true}
		btnSubscribe := markup.Data("", "")
		if status {
			btnSubscribe = markup.Data(unsubscribedText, unsubscribedEvent)
		} else {
			btnSubscribe = markup.Data(subscribedText, subscribedEvent)
		}
		markup.Inline(
			markup.Row(btnSubscribe), // одна кнопка в ряд
		)

		// ... сохраняем подписку в базу ...

		return c.Edit(
			"Вы успешно подписались на события!",
			markup,
		)
		// или c.Respond(&tele.CallbackResponse{Text: "Подписка оформлена!"})
	})

	// Обработка нажатия на кнопку
	bot.Handle(&telebot.Btn{Unique: unsubscribedEvent}, func(c telebot.Context) error {
		log.Printf("Unsubscribed| chat: %d, sender: %d, username: %s", c.Chat().ID, c.Sender().ID, c.Sender().Username)
		status, err := dao.service.Unsubscribed(c.Chat().ID)
		if err != nil {
			return err
		}

		markup := &telebot.ReplyMarkup{ResizeKeyboard: true}
		btnSubscribe := markup.Data("", "")
		if status {
			btnSubscribe = markup.Data(subscribedText, subscribedEvent)
		} else {
			btnSubscribe = markup.Data(unsubscribedText, unsubscribedEvent)
		}
		markup.Inline(
			markup.Row(btnSubscribe),
		)

		return c.Edit("Вы вы отменили подписку на события!",
			markup,
		)
	})

	bot.Start()

}
