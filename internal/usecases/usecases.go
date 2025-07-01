package usecases

import (
	"main-service/internal/sl"
	"math/rand"
)

type Services struct {
	UserService   *UserService
	AdminService *AdminService
	FamilyService *FamilyService
	InviteService *InviteService
}

func New(
	userRepo UserRepository,
	familyRepo FamilyRepository,
	inviteRepo InviteRepository,
	sl *sl.MyLogger,
) *Services {
	return &Services{
		UserService:   NewUserService(userRepo, sl),
		AdminService: NewAdminService(userRepo, familyRepo, sl),
		FamilyService: NewFamilyService(familyRepo, sl),
		InviteService: NewInviteService(inviteRepo, sl),
	}
}

const (
	codeLength = 6
)

var generateInviteCode = func() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, codeLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

