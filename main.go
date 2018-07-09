package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/halink0803/telegram-raffle-bot/common"
)

const (
	updateWalletAddress = "update_wallet_address"
)

var questions []Question
var chatGroup string

func readConfigFromFile(path string) (common.BotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return common.BotConfig{}, err
	}
	result := common.BotConfig{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func readQuestionsFromFile(path string) ([]Question, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []Question{}, err
	}
	result := []Question{}
	err = json.Unmarshal(data, &result)
	return result, err
}

//Bot the main bot
type Bot struct {
	bot     *tgbotapi.BotAPI
	storage *RaffleStorage
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	configPath := "config.json"
	botConfig, err := readConfigFromFile(configPath)
	if err != nil {
		log.Panic(err)
	}

	chatGroup = botConfig.ChatGroup

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

	storagePath := "/db/raffle.db"
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
		case "me":
			mybot.handleMe(update)
			break
		case "report":
			mybot.handleReport(update)
			break
		case "A":
			mybot.handleAnswer(update, 0)
			break
		case "B":
			mybot.handleAnswer(update, 1)
			break
		case "C":
			mybot.handleAnswer(update, 2)
			break
		case "D":
			mybot.handleAnswer(update, 3)
			break
		default:
			break
		}
	} else {
		userID := update.Message.From.ID
		currentAction, err := mybot.storage.GetCurrentAction(userID)
		if err != nil {
			log.Panic(err)
		}
		switch currentAction {
		case updateWalletAddress:
			mybot.handleUpdateWalletAddress(update)
			break
		}
	}
}

func replybuttons(numberOfOptions int) tgbotapi.ReplyKeyboardMarkup {
	firstRow := []tgbotapi.KeyboardButton{}
	secondRow := []tgbotapi.KeyboardButton{}

	firstChoice := tgbotapi.NewKeyboardButton("/A")
	firstRow = append(firstRow, firstChoice)

	secondChoice := tgbotapi.NewKeyboardButton("/B")
	firstRow = append(firstRow, secondChoice)

	if numberOfOptions == 2 {
		buttons := tgbotapi.NewReplyKeyboard(firstRow)
		return buttons
	}

	thirdChoice := tgbotapi.NewKeyboardButton("/C")
	secondRow = append(secondRow, thirdChoice)

	fourthChoice := tgbotapi.NewKeyboardButton("/D")
	secondRow = append(secondRow, fourthChoice)

	buttons := tgbotapi.NewReplyKeyboard(firstRow, secondRow)
	return buttons
}

func (mybot *Bot) finishAnswer(update tgbotapi.Update) {
	msgContent := fmt.Sprintf("Congratulation. You have answered all questions. Please leave your ETH address in the next messsage. We will contact and send prize if you win the raffle.")
	userID := update.Message.From.ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
	replyButtons := tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = replyButtons
	mybot.bot.Send(msg)
	err := mybot.storage.UpdateCurrentAction(userID, updateWalletAddress)
	if err != nil {
		log.Panic(err)
	}
}

func (mybot *Bot) nextQuestion(update tgbotapi.Update) {
	userID := update.Message.From.ID
	currentQuestion, err := mybot.storage.CurrentQuestion(userID)
	if err != nil {
		log.Panic(err)
	}
	currentQuestion++
	randQuestions, err := mybot.storage.GetRandQuestions(userID)
	if err != nil {
		log.Panic(err)
	}
	if currentQuestion > len(randQuestions)-1 {
		mybot.finishAnswer(update)
		return
	}
	question := questions[randQuestions[currentQuestion]]
	msgContent := fmt.Sprintf("%d. %s\n\n", currentQuestion+1, question.Question)
	options := []string{"/A", "/B", "/C", "/D"}
	for index, option := range question.Options {
		msgContent += fmt.Sprintf("%s %s\n", options[index], option)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
	buttons := replybuttons(len(question.Options))
	msg.ReplyMarkup = buttons
	mybot.bot.Send(msg)
	mybot.storage.UpdateCurrentQuestion(userID, currentQuestion)
}

//Generate 5 random questions in the list
func randQuestions() []int {
	// random a new sequence of question
	rand.Seed(time.Now().UnixNano())
	rands := rand.Perm(5)[:5]
	log.Printf("rands: %+v", rands)
	return rands
}

func (mybot *Bot) handleAnswer(update tgbotapi.Update, option int) {
	log.Printf("handle answer")
	userID := update.Message.From.ID
	currentQuestion, err := mybot.storage.CurrentQuestion(userID)
	log.Println(currentQuestion)
	if err != nil {
		log.Panic(err)
	}
	user, err := mybot.storage.GetUser(userID)
	if err != nil {
		log.Panic(err)
	}
	randQuestions := user.RandQuestions
	question := questions[randQuestions[currentQuestion]]
	if question.Answer == option {
		user.Score++
		if err := mybot.storage.UpdateUserScore(userID, user.Score); err != nil {
			log.Panic(err)
		}
		if err := mybot.storage.UpdateQuestionScore(question, 1); err != nil {
			log.Panic(err)
		}
	} else {
		if err := mybot.storage.UpdateQuestionScore(question, 0); err != nil {
			log.Panic(err)
		}
	}
	mybot.nextQuestion(update)
}

func (mybot *Bot) handleStart(update tgbotapi.Update) {
	userID := update.Message.From.ID
	randQuestions := randQuestions()
	user, err := mybot.storage.GetUser(userID)
	if err != nil {
		userName := fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)

		user = User{
			ID:              userID,
			Username:        userName,
			CurrentQuestion: -1,
			Valid:           true,
			Score:           0,
			RandQuestions:   randQuestions,
		}
		mybot.storage.AddUser(user)
	} else {
		if err := mybot.storage.UpdateRandQuestions(userID, randQuestions); err != nil {
			log.Panic(err)
		}
		if err := mybot.storage.UpdateCurrentQuestion(userID, -1); err != nil {
			log.Panic(err)
		}
		if err := mybot.storage.UpdateUserScore(userID, 0); err != nil {
			log.Panic(err)
		}
	}
	mybot.nextQuestion(update)
}

func (mybot *Bot) isAdminPrevilege(update tgbotapi.Update) bool {
	log.Print(chatGroup)
	chatConfig := tgbotapi.ChatConfigWithUser{
		SuperGroupUsername: "@" + chatGroup,
		UserID:             update.Message.From.ID,
	}
	chatMember, err := mybot.bot.GetChatMember(chatConfig)
	if err != nil {
		log.Panic(err)
	}
	return chatMember.IsAdministrator() || chatMember.IsCreator()
}

func (mybot *Bot) handleReport(update tgbotapi.Update) {
	//check for admin privilege
	if !mybot.isAdminPrevilege(update) {
		msgContent := "You are not privilege to run this command."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
		mybot.bot.Send(msg)
		return
	}
	//report following things:
	//how many overall correct answer
	totalUserAnswer, totalPoint, err := mybot.storage.UserStat()
	if err != nil {
		log.Panic(err)
	}
	//average score
	averageScore := 0.0
	if totalUserAnswer > 0 {
		averageScore = float64(totalPoint) / float64(totalUserAnswer)
	}
	//how many people answer correctly for each question
	msgContent := fmt.Sprintf("Total user answered questions: %d\n", totalUserAnswer)
	msgContent += fmt.Sprintf("Average score: %.2f\n", averageScore)

	//average score by question
	questions, err := mybot.storage.QuestionStat()
	if err != nil {
		log.Panic(err)
	}
	msgContent += fmt.Sprintf("Average score by question: \n")
	for _, question := range questions {
		average := 0.0
		if question.NumberUser > 0 {
			average = float64(question.Score) / float64(question.NumberUser)
		}
		msgContent += fmt.Sprintf("Question %d: %.2f\n", question.ID, average)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
	mybot.bot.Send(msg)
}

func (mybot *Bot) handleMe(update tgbotapi.Update) {
	userID := update.Message.From.ID
	user, err := mybot.storage.GetUser(userID)
	if err != nil {
		userName := fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
		user = User{
			ID:              userID,
			Username:        userName,
			CurrentQuestion: -1,
			Valid:           true,
			Score:           0,
		}
		err := mybot.storage.AddUser(user)
		if err != nil {
			log.Panic(err)
		}
	}
	msgContent := fmt.Sprintf("Your score: %d\n", user.Score)
	msgContent += fmt.Sprintf("Your wallet address: %s", user.WalletAddress)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
	mybot.bot.Send(msg)
}

func (mybot *Bot) handleUpdateWalletAddress(update tgbotapi.Update) {
	walletAddress := update.Message.Text
	userID := update.Message.From.ID
	err := mybot.storage.UpdateWalletAddress(userID, walletAddress)
	if err != nil {
		log.Panic(err)
	}

	//reset current action
	err = mybot.storage.UpdateCurrentAction(userID, "")
	if err != nil {
		log.Panic(err)
	}

	//response to user
	user, err := mybot.storage.GetUser(userID)
	if err != nil {
		log.Panic(err)
	}
	msgContent := fmt.Sprintf("Thank you. Your wallet address is: %s", user.WalletAddress)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
	mybot.bot.Send(msg)
}
