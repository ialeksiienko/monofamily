package usecase

import (
	"context"
	"monofamily/internal/entity"
	"monofamily/internal/errorsx"
	"monofamily/internal/service/userservice"
	"time"
)

type userService interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveUserToFamily(ctx context.Context, familyID int, userID int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUsersByFamilyID(ctx context.Context, familyID int) ([]entity.User, error)
	GetFamilyMembersInfo(ctx context.Context, family *entity.Family, userID int64) ([]userservice.MemberInfo, error)
	DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error
}

type adminService interface {
	DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error
}

type familyService interface {
	Create(ctx context.Context, familyName string, userID int64) (*entity.Family, error)
	GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error)
	GetFamilyByCode(ctx context.Context, code string) (*entity.Family, time.Time, error)
	GetFamilyByID(ctx context.Context, id int) (*entity.Family, error)
	CreateNewInviteCode(ctx context.Context, family *entity.Family, userID int64) (string, time.Time, error)
	DeleteFamily(ctx context.Context, familyID int) error
}

type UseCase struct {
	userService   userService
	adminService  adminService
	familyService familyService
}

func New(
	userService userService,
	adminService adminService,
	familyService familyService,
) *UseCase {
	return &UseCase{
		userService:   userService,
		adminService:  adminService,
		familyService: familyService,
	}
}

func (uc *UseCase) checkAdminPermission(family *entity.Family, userID int64) error {
	if userID != family.CreatedBy {
		return errorsx.NewError("no permission", errorsx.ErrCodeNoPermission, struct{}{})
	}
	return nil
}
