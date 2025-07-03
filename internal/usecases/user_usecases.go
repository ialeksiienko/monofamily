package usecases

import (
	"log/slog"
	"main-service/internal/entities"
	"main-service/internal/sl"
)

type UserProvider interface {
	GetAllUsersInFamily(familyID int) ([]entities.User, error)
	GetUserByID(id int64) (*entities.User, error)
}

type UserService struct {
	userSaver UserSaver
	userProvider UserProvider
	userDeletor UserDeletor
	sl *sl.MyLogger
}

func NewUserService(
	userSaver UserSaver,
	userProvider UserProvider,
	userDeletor UserDeletor,
	sl *sl.MyLogger,
) *UserService {
	return &UserService{
		userSaver: userSaver,
		userProvider: userProvider,
		userDeletor: userDeletor,
		sl: sl,
	}
}

func (s *UserService) Register(user *entities.User) (*entities.User, error) {
	return s.userSaver.SaveUser(user)
}

func (s *UserService) GetUserByID(id int64) (*entities.User, error) {
	return s.userProvider.GetUserByID(id)
}

func (s *UserService) GetUsers(familyID int) ([]entities.User, error) {
	return s.userProvider.GetAllUsersInFamily(familyID)
}

func (s *UserService) DeleteUserFromFamily(familyID int, userID int64) error {
	return s.userDeletor.DeleteUserFromFamily(familyID, userID)
}

type MemberInfo struct {
	ID        int64
	Username  string
	Firstname string
	IsAdmin bool
	IsCurrent bool
}

func (s *UserService) GetMembersInfo(family *entities.Family, userID int64) ([]MemberInfo, error) {
	users, err := s.userProvider.GetAllUsersInFamily(family.ID)
	if err != nil {
		s.sl.Error("failed to get all users in family", slog.String("family_name", family.Name),slog.String("err", err.Error()))
		return nil, err
	}

	if len(users) == 0 {
		return nil, &CustomError[struct{}]{
			Msg: "family has not members",
			Code: ErrCodeFamilyHasNoMembers,
		}
	}

	members := make([]MemberInfo, len(users))
	for i, user := range users {
		members[i] = MemberInfo{
			ID: user.ID,
			Username: user.Username,
			Firstname: user.Firstname,
			IsAdmin: family.CreatedBy == user.ID,
			IsCurrent: user.ID == userID,
		}
	}

	return members, nil
}

func (s *UserService) LeaveFamily(family *entities.Family, userID int64) error {
	if family.CreatedBy == userID {
		return &CustomError[struct{}]{
			Msg: "admin cannot leave family",
			Code: ErrCodeCannotRemoveSelf,
		}
	}

	err := s.userDeletor.DeleteUserFromFamily(family.ID, userID)
	if err != nil {
		s.sl.Error("failed to delete user from family", slog.Int("user_id", int(userID)), slog.Int("family_id", family.ID), slog.String("err", err.Error()))
		return err
	}

	return nil
}