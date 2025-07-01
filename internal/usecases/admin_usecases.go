package usecases

import (
	"log/slog"
	"main-service/internal/sessions"
	"main-service/internal/sl"
	"time"
)

type AdminService struct {
	userRepo UserRepository
	familyRepo FamilyRepository
	sl *sl.MyLogger
}

func NewAdminService(
	userRepo UserRepository,
	familyRepo FamilyRepository,
	sl *sl.MyLogger,
) *AdminService {
	return &AdminService{
		userRepo: userRepo,
		familyRepo: familyRepo,
		sl: sl,
	}
}

func(s *AdminService) RemoveMember(userID, memberID int64 ) error{
	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		return &CustomError[struct{}]{
			Msg: "user not in family",
			Code: ErrCodeUserNotInFamily,
		}
	}

	if userID != us.Family.CreatedBy {
		return &CustomError[struct{}]{
			Msg: "no permission",
			Code: ErrCodeNoPermission,
		}
	}

	if userID == memberID {
		return &CustomError[struct{}]{
			Msg: "cannot remove self",
			Code: ErrCodeCannotRemoveSelf,
		}
	}

	err := s.userRepo.DeleteUserFromFamily(us.Family.ID, memberID)
	if err != nil {
		s.sl.Error("unable to delete user from family", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *AdminService) DeleteFamily(userID int64) error {
	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		return &CustomError[struct{}]{
			Msg: "user not in family",
			Code: ErrCodeUserNotInFamily,
		}
	}

	if userID != us.Family.CreatedBy {
		return &CustomError[struct{}]{
			Msg: "no permission",
			Code: ErrCodeNoPermission,
		}
	}

	err := s.familyRepo.DeleteFamily(us.Family.ID)
	if err != nil {
		s.sl.Error("failed to delete family", slog.String("error", err.Error()))
		return err
	}

	sessions.DeleteUserState(userID)

	return nil
}

func (s *AdminService) CreateNewFamilyCode(userID int64) (string,time.Time, error) {
	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		return "", time.Time{}, &CustomError[struct{}]{
			Msg: "user not in family",
			Code: ErrCodeUserNotInFamily,
		}
	}

	if userID != us.Family.CreatedBy {
		return "", time.Time{}, &CustomError[struct{}]{
			Msg: "no permission",
			Code: ErrCodeNoPermission,
		}
	}

	code := generateInviteCode()

	expiresAt, err := s.familyRepo.SaveFamilyInviteCode(userID, us.Family.ID, code)
	if err != nil {
		s.sl.Error("failed to save family invite code", slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	return code, expiresAt, nil
}