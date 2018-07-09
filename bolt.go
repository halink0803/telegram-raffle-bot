package main

import (
	"github.com/asdine/storm"
)

//User represent a user
type User struct {
	ID              int `storm:"id"`
	Username        string
	Name            string
	CurrentQuestion int `storm:"index"`
	Valid           bool
	Score           int
	RandQuestions   []int
	CurrentAction   string
	WalletAddress   string
}

//Question questions list
type Question struct {
	ID         int    `storm:"id"`
	Question   string `json:"question"`
	Options    []string
	Answer     int
	Score      int
	NumberUser int
}

//RaffleStorage represent bot storage
type RaffleStorage struct {
	db *storm.DB
}

//NewBoltStorage represent storage for the bot
func NewBoltStorage(path string) (*RaffleStorage, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	storage := &RaffleStorage{
		db: db,
	}
	return storage, nil
}

//CurrentQuestion return current question of a user
func (r *RaffleStorage) CurrentQuestion(userID int) (int, error) {
	var user User
	err := r.db.One("ID", userID, &user)
	return user.CurrentQuestion, err
}

//UpdateCurrentQuestion update user current question
func (r *RaffleStorage) UpdateCurrentQuestion(userID, currentQuestion int) error {
	return r.db.UpdateField(&User{ID: userID}, "CurrentQuestion", currentQuestion)
}

//GetUser return one user instant with userID
func (r *RaffleStorage) GetUser(userID int) (User, error) {
	var user User
	err := r.db.One("ID", userID, &user)
	return user, err
}

//AddUser and one user to storage
func (r *RaffleStorage) AddUser(user User) error {
	return r.db.Save(&user)
}

//UpdateUserScore update user score
func (r *RaffleStorage) UpdateUserScore(userID, score int) error {
	return r.db.UpdateField(&User{ID: userID}, "Score", score)
}

//UpdateRandQuestions update random questions for an user
func (r *RaffleStorage) UpdateRandQuestions(userID int, randQuestions []int) error {
	return r.db.UpdateField(&User{ID: userID}, "RandQuestions", randQuestions)
}

//GetRandQuestions return an array of rand questionsfor an user
func (r *RaffleStorage) GetRandQuestions(userID int) ([]int, error) {
	var user User
	err := r.db.One("ID", userID, &user)
	return user.RandQuestions, err
}

//UpdateCurrentAction update current action of an user
func (r *RaffleStorage) UpdateCurrentAction(userID int, action string) error {
	return r.db.UpdateField(&User{ID: userID}, "CurrentAction", action)
}

//GetCurrentAction get current action which one user is currently doing
func (r *RaffleStorage) GetCurrentAction(userID int) (string, error) {
	var user User
	err := r.db.One("ID", userID, &user)
	return user.CurrentAction, err
}

//UpdateWalletAddress update wallet address of an user
func (r *RaffleStorage) UpdateWalletAddress(userID int, walletAddress string) error {
	return r.db.UpdateField(&User{ID: userID}, "WalletAddress", walletAddress)
}

//Report return a report
func (r *RaffleStorage) Report() error {
	return nil
}

//UserStat return total user score and total point
func (r *RaffleStorage) UserStat() (int, int, error) {
	var users []User
	if err := r.db.All(&users); err != nil {
		return 0, 0, err
	}
	totalUserAnswered := len(users)
	totalPoint := 0
	for _, user := range users {
		totalPoint += user.Score
	}
	return totalUserAnswered, totalPoint, nil
}

//UpdateQuestionScore update question score
func (r *RaffleStorage) UpdateQuestionScore(question Question, score int) error {
	var q Question
	err := r.db.One("ID", question.ID, &q)
	if err != nil {
		question.Score += score
		question.NumberUser = 1
		err = r.db.Save(&question)
		if err != nil {
			return err
		}
		return nil
	}
	if err = r.db.UpdateField(&q, "Score", q.Score+score); err != nil {
		return err
	}
	err = r.db.UpdateField(&q, "NumberUser", q.NumberUser+1)
	return err
}

//QuestionStat return stat of questions
func (r *RaffleStorage) QuestionStat() ([]Question, error) {
	var questions []Question
	err := r.db.AllByIndex("ID", &questions)
	return questions, err
}
