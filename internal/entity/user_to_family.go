package entity

type UserToFamily struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	FamilyID int `json:"family_id"`
}
