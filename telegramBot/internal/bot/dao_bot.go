package bot

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"telegramBot/internal/config"
	customError "telegramBot/internal/error"
	"telegramBot/internal/service"
	"time"

	"gopkg.in/telebot.v3"
)

type State string

const (
	StateMain        State = "state_main"
	StateSpaces      State = "state_spaces"
	StateProposals   State = "state_proposals"
	StateVotes       State = "state_votes"
	StateVotesDelete State = "state_votes_delete"
)

const (
	spacesText    = "Spaces"
	proposalsText = "Proposals"
	votesText     = "Votes"

	cancelText = "⬅ В меню"

	subscribedText   = "Подписаться"
	unsubscribedText = "Отписаться"

	deletedVotes = "Удалить votes"
)

const (
	spacesEvent    = "spaces_event"
	proposalsEvent = "proposals_event"
	votesEvent     = "votes_event"

	cancelEvent      = "cancel_event"
	votesDeleteEvent = "votes_delete_event"

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

var (
	fsm = map[int64]State{}
)

func (dao DaoBot) GetState(id int64) State {
	if s, ok := fsm[id]; ok {
		return s
	}
	return StateMain
}

func (dao DaoBot) SetState(id int64, s State) {
	fsm[id] = s
}

func (dao DaoBot) MainMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(
			m.Data("Spaces", spacesEvent),
			m.Data("Proposals", proposalsEvent),
			m.Data("Votes", votesEvent),
		),
	)
	return m
}

func (dao DaoBot) SubscribeMenu(subscribed bool) *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}

	title := subscribedText
	event := subscribedEvent
	if subscribed {
		title = unsubscribedText
		event = unsubscribedEvent
	}

	m.Inline(
		m.Row(
			m.Data(cancelText, cancelEvent),
			m.Data(title, event),
		),
	)
	return m
}

func (dao DaoBot) VotesMenu(deleteButton bool) *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}

	rowMenu := m.Row(m.Data(cancelText, cancelEvent))
	if deleteButton {
		rowMenu = m.Row(
			m.Data(cancelText, cancelEvent),
			m.Data(deletedVotes, votesDeleteEvent),
		)
	}

	m.Inline(
		rowMenu,
	)
	return m
}

func (dao DaoBot) Callbacks(c telebot.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data[1:]

	switch data {
	case spacesEvent:
		dao.SetState(userID, StateSpaces)
		status, err := dao.service.StatusSubscribedSpaces(userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err := dao.service.NewUser(userID)
				if err != nil {
					return err
				}
			}
			return err
		}
		var statusBool = status == 1
		return c.Edit("Spaces", dao.SubscribeMenu(statusBool))

	case proposalsEvent:
		dao.SetState(userID, StateProposals)
		status, err := dao.service.StatusSubscribedProposals(userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err := dao.service.NewUser(userID)
				if err != nil {
					return err
				}
			}
			return err
		}
		var statusBool = status == 1
		return c.Edit("Proposals", dao.SubscribeMenu(statusBool))

	case votesEvent:
		dao.SetState(userID, StateVotes)
		votes, err := dao.service.GetVotesByUser(userID)

		if err != nil {
			if errors.Is(err, customError.ErrorNotFound) {
				return c.Edit("Введите id votes", dao.VotesMenu(false))
			}
			return err
		}
		outMessages := ""
		for _, vote := range votes {
			outMessages += vote + "\n"
		}
		outMessages += "Введите id votes"

		return c.Edit(outMessages, dao.VotesMenu(true))

	case subscribedEvent:
		fmt.Println("subscribedEvent")
		switch dao.GetState(userID) {
		case StateSpaces:
			status, err := dao.service.SubscribedSpaces(userID)
			if err != nil {
				return err
			}
			return c.Edit("Spaces", dao.SubscribeMenu(status))

		case StateProposals:
			status, err := dao.service.SubscribedProposals(userID)
			if err != nil {
				return err
			}
			return c.Edit("Proposals", dao.SubscribeMenu(status))
		}

	case unsubscribedEvent:
		fmt.Println("unsubscribedEvent")
		switch dao.GetState(userID) {
		case StateSpaces:
			status, err := dao.service.UnsubscribedSpaces(userID)
			if err != nil {
				return err
			}
			return c.Edit("Spaces", dao.SubscribeMenu(!status))

		case StateProposals:
			status, err := dao.service.UnsubscribedProposals(userID)
			if err != nil {
				return err
			}
			return c.Edit("Proposals", dao.SubscribeMenu(!status))
		}

	case cancelEvent:
		dao.SetState(userID, StateMain)
		return c.Edit("Главное меню", dao.MainMenu())

	case votesDeleteEvent:
		votes, err := dao.service.GetVotesByUser(userID)
		if err != nil {
			if errors.Is(err, customError.ErrorNotFound) {
				dao.SetState(userID, StateVotes)

				return c.Edit("Id votes не найдены. Введите id votes", dao.VotesMenu(false))
			}
			return err
		}

		dao.SetState(userID, StateVotesDelete)
		outMessages := ""
		for _, vote := range votes {
			outMessages += vote + "\n"
		}
		outMessages += "Введите id votes для удаления"
		return c.Edit(outMessages, dao.VotesMenu(true))

	}

	return nil
}

func (dao DaoBot) StartHandle(c telebot.Context) error {
	log.Printf("Start| chat: %d, sender: %d, username: %s", c.Chat().ID, c.Sender().ID, c.Sender().Username)
	if err := dao.service.NewUser(c.Chat().ID); err != nil {
		log.Println("Start: ", err)
	}
	// Создаём клавиатуру
	markup := &telebot.ReplyMarkup{ResizeKeyboard: true}

	btnSpacesSubscribe := markup.Data(spacesText, spacesEvent)
	btnProposalsSubscribe := markup.Data(proposalsText, proposalsEvent)
	btnVotesSubscribe := markup.Data(votesText, votesEvent)

	markup.Inline(
		markup.Row(btnSpacesSubscribe, btnProposalsSubscribe, btnVotesSubscribe),
	)
	return c.Send(
		"Какие уведомления хочешь получать?",
		markup,
	)
}

func (dao DaoBot) Messages(c telebot.Context) error {
	userID := c.Sender().ID
	votesID := c.Message().Text

	fmt.Println("Messages ", userID, votesID, dao.GetState(userID))

	switch dao.GetState(userID) {
	case StateVotes:
		dao.SetState(userID, StateVotes)

		status, err := dao.service.CreateVotesId(userID, votesID)
		if err != nil {
			return err
		}

		id := strings.TrimSpace(c.Text())
		if id == "" {
			return c.Send("Введите корректный id.\nВведите id votes.")
		}

		outMessage := ""
		if status {
			outMessage = fmt.Sprintf("ID %s добавлен.\nВведите id votes еще.", id)
		} else {
			outMessage = fmt.Sprintf("Ошибка. Попробуйте еще")
		}

		if !status {
			_, err := dao.service.GetVotesByUser(userID)
			if err != nil {
				if errors.Is(err, customError.ErrorNotFound) {
					return c.Edit(outMessage, dao.VotesMenu(false))
				}
				return err
			}
		}
		return c.Send(
			outMessage,
			dao.VotesMenu(true),
		)
	case StateVotesDelete:
		dao.SetState(userID, StateVotesDelete)

		status, err := dao.service.DropVotesId(userID, votesID)
		if err != nil {
			return err
		}

		id := strings.TrimSpace(c.Text())
		if id == "" {
			return c.Send("Введите корректный id.\nВведите id votes.")
		}

		outMessage := ""
		if status {
			outMessage = fmt.Sprintf("ID %s удален.\nВведите id votes еще.", id)
		} else {
			outMessage = fmt.Sprintf("Ошибка. Попробуйте еще. \nВведите id votes еще.")
		}

		if !status {
			_, err := dao.service.GetVotesByUser(userID)
			if err != nil {
				if errors.Is(err, customError.ErrorNotFound) {
					return c.Edit(outMessage, dao.VotesMenu(false))
				}
				return err
			}
		}
		return c.Send(
			outMessage,
			dao.VotesMenu(true),
		)
	}

	return nil
}

func (dao DaoBot) InitRoutes(bot *telebot.Bot) {

	bot.Handle("/start", dao.StartHandle)
	bot.Handle(telebot.OnCallback, dao.Callbacks)
	bot.Handle(telebot.OnText, dao.Messages)
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

	dao.InitRoutes(bot)

	bot.Start()

}
