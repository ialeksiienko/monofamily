package usecases

import (
	"crypto/rand"
	"encoding/base32"
	"main-service/internal/sl"
	"strings"
)

type Services struct {
	UserService   *UserService
	AdminService *AdminService
	FamilyService *FamilyService
	FamilyInviteCodeService *FamilyInviteCodeService
}

func New(
	userRepo UserRepository,
	familyRepo FamilyRepository,
	familyInviteCodeRepo FamilyInviteCodeRepository,
	sl *sl.MyLogger,
) *Services {
	return &Services{
		UserService:   NewUserService(userRepo, userRepo, userRepo, sl),
		AdminService: NewAdminService(userRepo, familyRepo, familyInviteCodeRepo, sl),
		FamilyService: NewFamilyService(userRepo, familyRepo, familyRepo, familyInviteCodeRepo, sl),
		FamilyInviteCodeService: NewFamilyInviteCodeService(familyInviteCodeRepo, sl),
	}
}

const (
	codeLength = 6
)

var generateInviteCode = func() (string, error) {
	b := make([]byte, 5)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	code := base32.StdEncoding.EncodeToString(b)
	code = strings.ToUpper(code)
	return code[:codeLength], nil
}

