package models

type UsersToFamily struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`
	FamilyID uint64 `json:"family_id"`
}
