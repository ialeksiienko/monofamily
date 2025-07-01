package usecases

import (
	"log/slog"
	"main-service/internal/sl"
)


type InviteService struct {
	inviteRepo InviteRepository
	sl *sl.MyLogger
}

func NewInviteService(inviteRepo InviteRepository, sl *sl.MyLogger) *InviteService {
	return &InviteService{inviteRepo: inviteRepo, sl: sl}
}

func (s *InviteService) ClearInviteCodes() error {
	err := s.inviteRepo.ClearInviteCodes()
	if err != nil {
		s.sl.Error("failed to clear invite codes", slog.String("error", err.Error()))
		return err
	}

	return nil
}
