package usecases

import (
	"log/slog"
	"main-service/internal/entities"
	"main-service/internal/sessions"
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

func (s *UserService) GetMembersInfo(userID int64) ([]MemberInfo, error) {
	us, exists := sessions.GetUserState(userID)
	if !exists {
		return nil, &CustomError[struct{}]{
			Msg: "user not in family",
			Code: ErrCodeUserNotInFamily,
		}
	}

	users, err := s.userRepo.GetAllUsersInFamily(us.Family.ID)
	if err != nil {
		s.sl.Error("failed to get all users in family", slog.String("error", err.Error()))
		return nil, err
	}

	members := make([]MemberInfo, len(users))
	for i, user := range users {
		members[i] = MemberInfo{
			ID: user.ID,
			Username: user.Username,
			Firstname: user.Firstname,
			IsAdmin: us.Family.CreatedBy == userID,
			IsCurrent: user.ID == userID,
		}
	}

	return members, nil
}

func (s *UserService) LeaveFamily(userID int64) error {
	us, exists := sessions.GetUserState(userID)
	if !exists {
		return &CustomError[struct{}]{
			Msg: "user not in family",
			Code: ErrCodeUserNotInFamily,
		}
	}

	if us.Family.CreatedBy == userID {
		return &CustomError[struct{}]{
			Msg: "admin cannot remove self from family",
			Code: ErrCodeCannotRemoveSelf,
		}
	}

	err := s.userRepo.DeleteUserFromFamily(us.Family.ID, userID)
	if err != nil {
		return err
	}

	sessions.DeleteUserState(userID)

	return nil
}