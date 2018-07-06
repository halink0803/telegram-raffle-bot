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
func (r *RaffleStorage) Report() {

}
