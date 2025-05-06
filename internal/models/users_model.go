package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	JoinedAt  time.Time `json:"joined_at"`
	//Token    string `json:"token"`
	//CardType string `json:"card_type"`
	//Currency string `json:"currency"`
	//Balance  string `json:"balance"`
}
