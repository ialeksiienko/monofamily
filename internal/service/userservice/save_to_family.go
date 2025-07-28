package userservice

import (
	"context"
	"log/slog"
)

func (s *UserService) SaveUserToFamily(ctx context.Context, familyID int, userID int64) error {
	saveErr := s.userSaver.SaveUserToFamily(ctx, familyID, userID)
	if saveErr != nil {
		s.sl.Error("unable to save user to family", slog.Int("user_id", int(userID)), slog.Int("family_id", familyID), slog.String("err", saveErr.Error()))
		return saveErr
	}
	return nil
}
