package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/halink0803/telegram-raffle-bot/common"
)

var questions []common.Questions

func readConfigFromFile(path string) (common.BotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return common.BotConfig{}, err
	}
	result := common.BotConfig{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func readQuestionsFromFile(path string) ([]common.Questions, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []common.Questions{}, err
	}
	result := []common.Questions{}
	err = json.Unmarshal(data, &result)
	return result, err
}

//Bot the main bot
type Bot struct {
	bot     *tgbotapi.BotAPI
	storage *RaffleStorage
}

func main() {
	configPath := "config.json"
	botConfig, err := readConfigFromFile(configPath)
	if err != nil {
		log.Panic(err)
	}

	// init bot
	bot, err := tgbotapi.NewBotAPI(botConfig.BotKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	questionPath := "questions.json"
	questions, err = readQuestionsFromFile(questionPath)
	if err != nil {
		log.Panic(err)
	}

	storagePath := "raffle.db"
	storage, err := NewBoltStorage(storagePath)
	if err != nil {
		log.Panic(err)
	}

	mybot := Bot{
		bot:     bot,
		storage: storage,
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := mybot.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("Updates:[%s] %s", update.Message.From.UserName, update.Message.Command())

		mybot.handle(update)

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// mybot.bot.Send(msg)
	}
}

func (mybot *Bot) handle(update tgbotapi.Update) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			mybot.handleStart(update)
			break
		case "report":
			mybot.handleReport(update)
		default:
			break
		}
	}
}

func replybuttons() tgbotapi.ReplyKeyboardMarkup {
	firstRow := []tgbotapi.KeyboardButton{}
	secondRow := []tgbotapi.KeyboardButton{}

	firstChoice := tgbotapi.NewKeyboardButton("A")
	firstRow = append(firstRow, firstChoice)

	secondChoice := tgbotapi.NewKeyboardButton("B")
	firstRow = append(firstRow, secondChoice)

	thirdChoice := tgbotapi.NewKeyboardButton("C")
	secondRow = append(secondRow, thirdChoice)

	fourthChoice := tgbotapi.NewKeyboardButton("D")
	secondRow = append(secondRow, fourthChoice)

	buttons := tgbotapi.NewReplyKeyboard(firstRow, secondRow)
	return buttons
}

//TODO: send next question
func (mybot *Bot) nextQuestion(update tgbotapi.Update) {
	// userID := update.Message.From.ID
	// currentQuestion := mybot.storage.CurrentQuestion()
}

func randQuestions() {

}

func (mybot *Bot) handleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hihi")
	buttons := replybuttons()
	msg.ReplyMarkup = buttons
	mybot.bot.Send(msg)
}

func (mybot *Bot) handleReport(update tgbotapi.Update) {
	//report following things:
	//how many overall correct answer

	//average score

	//how many people answer correctly for each question
}
