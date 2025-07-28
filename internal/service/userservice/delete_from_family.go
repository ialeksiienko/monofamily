package userservice

import (
	"context"
	"log/slog"
)

func (s *UserService) DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error {
	err := s.userDeletor.DeleteUserFromFamily(ctx, familyID, userID)
	if err != nil {
		s.sl.Error("unable to delete member from family", slog.Int("member_id", int(userID)), slog.Int("family_id", familyID), slog.String("error", err.Error()))
		return err
	}

	return nil
}
