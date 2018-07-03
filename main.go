package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/halink0803/telegram-raffle-bot/common"
)

var questions common.Questions

func readConfigFromFile(path string) (common.BotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return common.BotConfig{}, err
	}
	result := common.BotConfig{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func readQuestionsFromFile(path string) (common.Questions, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return common.Questions{}, err
	}
	result := common.Questions{}
	err = json.Unmarshal(data, &result)
	return result, err
}

//Bot the main bot
type Bot struct {
	bot *tgbotapi.BotAPI
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

	questionPath := "question.json"
	questions, err = readQuestionsFromFile(questionPath)
	if err != nil {
		log.Panic(err)
	}

	mybot := Bot{
		bot: bot,
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
	replyRow := []tgbotapi.KeyboardButton{}
	//like button
	likeKeyboardButton := tgbotapi.NewKeyboardButton("A")
	replyRow = append(replyRow, likeKeyboardButton)

	//unlike button
	unlikeKeyboardButton := tgbotapi.NewKeyboardButton("B")
	replyRow = append(replyRow, unlikeKeyboardButton)

	//download button
	downloadKeyboardButton := tgbotapi.NewKeyboardButton("C")
	replyRow = append(replyRow, downloadKeyboardButton)

	buttons := tgbotapi.NewReplyKeyboard(replyRow)
	return buttons
}

func (mybot *Bot) handleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hihi")
	buttons := replybuttons()
	msg.ReplyMarkup = buttons
	mybot.bot.Send(msg)
}

func (mybot *Bot) handleReport(update tgbotapi.Update) {
}
