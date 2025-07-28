package userservice

import (
	"context"
	"monofamily/internal/entity"
	"monofamily/internal/pkg/sl"
)

type UserServiceIface interface {
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveUserToFamily(ctx context.Context, familyID int, userID int64) error
	GetAllUsersInFamily(ctx context.Context, familyID int) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error
}

type userSaver interface {
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveUserToFamily(ctx context.Context, familyID int, userID int64) error
}

type userProvider interface {
	GetAllUsersInFamily(ctx context.Context, familyID int) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
}

type userDeletor interface {
	DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error
}

type UserService struct {
	userSaver    userSaver
	userProvider userProvider
	userDeletor  userDeletor
	sl           sl.Logger
}

func New(
	userIface UserServiceIface,
	sl sl.Logger,
) *UserService {
	return &UserService{
		userSaver:    userIface,
		userProvider: userIface,
		userDeletor:  userIface,
		sl:           sl,
	}
}
