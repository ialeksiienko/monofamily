package entity

import "time"

type UserBankToken struct {
	ID        int
	UserID    int64
	FamilyID  int
	Token     string
	CreatedAt time.Time
}
