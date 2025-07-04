package usecases

import (
	"main-service/internal/entities"
	"time"
)

type UserRepository interface {
	UserSaver
	UserProvider
	UserDeletor
}

type FamilyInviteCodeRepository interface {
	FamilyInviteCodeSaver
	FamilyInviteCodeCleaner
}

type FamilyRepository interface {
	FamilyCreator
	FamilyProvider
	FamilyDeletor
}

type UserSaver interface {
	SaveUser(user *entities.User) (*entities.User, error)
	SaveUserToFamily(familyID int, userID int64) error
}

type UserDeletor interface {
	DeleteUserFromFamily(familyID int, userID int64) error
}

type FamilyInviteCodeSaver interface {
	SaveFamilyInviteCode(userId int64, familyId int, code string) (time.Time, error)
}

type FamilyDeletor interface {
	DeleteFamily(familyID int) error
}
