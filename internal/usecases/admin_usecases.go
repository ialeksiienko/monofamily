package usecases

import (
	"log/slog"
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

func(s *AdminService) RemoveMember(familyID int, memberID int64 ) error{
	err := s.userRepo.DeleteUserFromFamily(familyID, memberID)
	if err != nil {
		s.sl.Error("unable to delete member from family", slog.Int("member_id", int(memberID)), slog.Int("family_id", familyID), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *AdminService) DeleteFamily(familyID int) error {
	err := s.familyRepo.DeleteFamily(familyID)
	if err != nil {
		s.sl.Error("failed to delete family", slog.Int("family_id", familyID), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *AdminService) CreateNewFamilyCode(familyID int, userID int64) (string,time.Time, error) {
	code := generateInviteCode()

	expiresAt, err := s.familyRepo.SaveFamilyInviteCode(userID, familyID, code)
	if err != nil {
		s.sl.Error("failed to save family invite code", slog.Int("created_by", int(userID)), slog.Int("family_id", familyID), slog.String("code", code), slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	return code, expiresAt, nil
}