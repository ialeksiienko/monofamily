package usecases

import (
	"log/slog"
	"main-service/internal/entities"
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

func(s *AdminService) RemoveMember(family *entities.Family, userID int64, memberID int64 ) error{
	if userID != family.CreatedBy {
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

	err := s.userRepo.DeleteUserFromFamily(family.ID, memberID)
	if err != nil {
		s.sl.Error("unable to delete member from family", slog.Int("member_id", int(memberID)), slog.Int("family_id", family.ID), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *AdminService) DeleteFamily(family *entities.Family, userID int64) error {
	if userID != family.CreatedBy {
		return &CustomError[struct{}]{
			Msg: "no permission",
			Code: ErrCodeNoPermission,
		}
	}

	err := s.familyRepo.DeleteFamily(family.ID)
	if err != nil {
		s.sl.Error("failed to delete family", slog.Int("family_id", family.ID), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *AdminService) CreateNewFamilyCode(family *entities.Family, userID int64) (string,time.Time, error) {
	if userID != family.CreatedBy {
		return "", time.Time{}, &CustomError[struct{}]{
			Msg: "no permission",
			Code: ErrCodeNoPermission,
		}
	}

	code := generateInviteCode()

	expiresAt, err := s.familyRepo.SaveFamilyInviteCode(userID, family.ID, code)
	if err != nil {
		s.sl.Error("failed to save family invite code", slog.Int("created_by", int(userID)), slog.Int("family_id", family.ID), slog.String("code", code), slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	return code, expiresAt, nil
}