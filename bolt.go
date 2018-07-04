package main

import (
	"github.com/asdine/storm"
)

//User represent a user
type User struct {
	ID              int `storm:"id,increment"`
	UserID          int `storm:"index"`
	Username        string
	Name            string
	CurrentQuestion int `storm:"index"`
	Valid           bool
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
	return nil
}

//Report return a report
func (r *RaffleStorage) Report() {

}
