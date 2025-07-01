package usecases

import (
	"main-service/internal/entities"
	"time"
)

type UserRepository interface {
	SaveUser(user *entities.User) (*entities.User, error)
	GetAllUsersInFamily(familyID int) ([]entities.User, error)
	DeleteUserFromFamily( familyID int, userID int64) error
}

type FamilyRepository interface {
	CreateFamily(inp *entities.Family) (*entities.Family, error)
	SaveUserToFamily(familyID int, userID int64) error
	SaveFamilyInviteCode(userId int64, familyId int, code string) (time.Time, error)
	GetFamiliesByUserID(userID int64) ([]entities.Family, error)
	GetFamilyByCode(code string) (*entities.Family, time.Time, error)
	GetFamilyBy(by string, value any) (*entities.Family, error)
	DeleteFamily(familyID int) error
}

type InviteRepository interface {
	ClearInviteCodes() error
}
