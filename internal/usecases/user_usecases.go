package usecases

import (
	"log/slog"
	"main-service/internal/entities"
	"main-service/internal/sl"
)

type UserService struct {
	userRepo UserRepository
	sl *sl.MyLogger
}

func NewUserService(userRepo UserRepository, sl *sl.MyLogger) *UserService {
	return &UserService{userRepo: userRepo, sl: sl}
}

func (s *UserService) Register(user *entities.User) (*entities.User, error) {
	return s.userRepo.SaveUser(user)
}

func (s *UserService) GetUsers(familyID int) ([]entities.User, error) {
	return s.userRepo.GetAllUsersInFamily(familyID)
}

func (s *UserService) DeleteUserFromFamily(familyID int, userID int64) error {
	return s.userRepo.DeleteUserFromFamily(familyID, userID)
}

type MemberInfo struct {
	ID        int64
	Username  string
	Firstname string
	IsAdmin bool
	IsCurrent bool
}

func (s *UserService) GetMembersInfo(family *entities.Family, userID int64) ([]MemberInfo, error) {
	users, err := s.userRepo.GetAllUsersInFamily(family.ID)
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
			IsAdmin: family.CreatedBy == userID,
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

	err := s.userRepo.DeleteUserFromFamily(family.ID, userID)
	if err != nil {
		s.sl.Error("failed to delete user from family", slog.Int("user_id", int(userID)), slog.Int("family_id", family.ID), slog.String("err", err.Error()))
		return err
	}

	return nil
}