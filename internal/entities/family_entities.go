package entities

import "time"

type Family struct {
	ID        int    `json:"id"`
	CreatedBy int64  `json:"created_by"`
	Name      string `json:"name"`
}

type FamilyInvites struct {
	ID        int       `json:"id"`
	FamilyID  int       `json:"family_id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int64     `json:"created_by"`
	ExpiresAt time.Time `json:"expires_at"`
}
