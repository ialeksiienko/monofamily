package entities

type UsersToFamily struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	FamilyID int `json:"family_id"`
}
