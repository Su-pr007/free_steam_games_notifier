package tgBot

import (
	"fmt"
	"log"
	"log/slog"
	"steamGamesSales/internal/entity"
	repository "steamGamesSales/internal/repository/postgres"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotData struct {
	log               *slog.Logger
	token             string
	repo              *repository.Repository
	newFreeGamesEvent chan []entity.Game
	bot               *tgbotapi.BotAPI
}

func mockBotAPI() *tgbotapi.BotAPI {
	return &tgbotapi.BotAPI{
		Token:  "",
		Debug:  false,
		Buffer: 0,
		Self:   tgbotapi.User{},
		Client: nil,
	}
}

func NewBot(log *slog.Logger, token string, repo *repository.Repository, newFreeGamesEvent chan []entity.Game) *BotData {
	return &BotData{
		log,
		token,
		repo,
		newFreeGamesEvent,
		mockBotAPI(),
	}
}
func (botData BotData) InitBot() *BotData {
	botObject, err := tgbotapi.NewBotAPI(botData.token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authorized on account %s\n", botObject.Self.UserName)

	botData.bot = botObject

	return &botData
}

func (botData BotData) StartListening() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botData.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("Text: %s\n", update.Message.Text)
		if update.Message.Text == "/start" {
			startHandler(update, botData)
		}
		// TODO: stop?
		if update.Message.Text == "/subscribe" {
			subscribeHandler(update, botData)
		}
		if update.Message.Text == "/unsubscribe" {
			unsubscribeHandler(update, botData)
		}
	}
}

func (botData BotData) StartWaitingForGames() {
	for games := range botData.newFreeGamesEvent {
		notificationText := formatText(games)
		err := botData.SendTextToAll(notificationText)
		if err != nil {
			return
		}
	}
}

func formatText(gamesList []entity.Game) string {
	result := "Бесплатные игры: \n"

	for _, game := range gamesList {
		result += fmt.Sprintf("- <b>%s %s</b>\n", game.Name, game.Link)
	}

	return result
}

func startHandler(update tgbotapi.Update, b BotData) bool {
	chat, _ := b.repo.Chat.Get(strconv.Itoa(int(update.Message.From.ID)))

	if chat != nil {
		_ = b.SendTextTo("Уже здоровались", update.Message.From.ID)

		return true
	}

	user := entity.Chat{
		Id:           update.Message.From.ID,
		Username:     update.Message.From.FirstName,
		Name:         update.Message.From.UserName,
		Subscribed:   false,
		LanguageCode: update.Message.From.LanguageCode,
		DateAdd:      update.Message.Date,
	}

	err := b.repo.Chat.Create(&user)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		return false
	}

	fmt.Printf("[%s] %s \n chatId: %d \n", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

	err = b.SendTextTo("Привееет", update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		return false
	}

	return true
}

func subscribeHandler(update tgbotapi.Update, b BotData) bool {
	chat, err := b.repo.Chat.Get(strconv.Itoa(int(update.Message.From.ID)))
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		return false
	}

	if chat.Subscribed == true {
		_ = b.SendTextTo("Вы уже подписаны на уведомления", update.Message.From.ID)

		return true
	}

	chat.Subscribed = true
	err = b.repo.Chat.Update(chat)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		return false
	}

	err = b.SendTextTo("Вы успешно подписаны на уведомления о раздачах игр", update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		return false
	}

	return true
}

func unsubscribeHandler(update tgbotapi.Update, b BotData) bool {
	chat, err := b.repo.Chat.Get(strconv.Itoa(int(update.Message.From.ID)))
	if err != nil {
		return false
	}

	if chat.Subscribed == false {
		err = b.SendTextTo("Вы не подписаны на уведомления", update.Message.From.ID)

		return true
	}

	chat.Subscribed = false
	err = b.repo.Chat.Update(chat)
	if err != nil {
		return false
	}

	err = b.SendTextTo("Вы отписаны от уведомлений", update.Message.From.ID)
	if err != nil {
		return false
	}

	return true
}

func (botData BotData) SendTextTo(text string, userId int64) error {
	msg := tgbotapi.NewMessage(userId, text)

	msg.ParseMode = "HTML"

	_, err := botData.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (botData BotData) SendTextToAll(text string) error {
	chats, err := botData.repo.Chat.GetSubscribed()
	if err != nil {
		panic(err)
	}

	for _, chat := range chats {
		msg := tgbotapi.NewMessage(chat.Id, text)

		msg.ParseMode = "HTML"

		_, err := botData.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

type ResponseData struct {
	Status string `json:"status"`
	Result struct {
		ImageID     string `json:"imageId"`
		MarketName  string `json:"marketName"`
		InspectLink string `json:"inspectLink"`
		State       string `json:"state"`
		Meta        struct {
			Images []struct {
				Slot int    `json:"slot"`
				Name string `json:"name"`
				Wear int    `json:"wear"`
			} `json:"5"`
		} `json:"meta"`
	} `json:"result"`
}
