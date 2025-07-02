package usecases

import (
	"log/slog"
	"main-service/internal/sl"
)

type FamilyInviteCodeCleaner interface {
	ClearInviteCodes() error
}

type FamilyInviteCodeService struct {
	familyInviteCodeCleaner FamilyInviteCodeCleaner
	sl *sl.MyLogger
}

func NewFamilyInviteCodeService(familyInviteCodeCleaner FamilyInviteCodeCleaner, sl *sl.MyLogger) *FamilyInviteCodeService {
	return &FamilyInviteCodeService{familyInviteCodeCleaner: familyInviteCodeCleaner, sl: sl}
}

func (s *FamilyInviteCodeService) ClearInviteCodes() error {
	err := s.familyInviteCodeCleaner.ClearInviteCodes()
	if err != nil {
		s.sl.Error("failed to clear invite codes", slog.String("error", err.Error()))
		return err
	}

	return nil
}
