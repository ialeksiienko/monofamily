package userrepo

import (
	"context"
	"log/slog"
)

func (ur *UserRepository) DeleteUserFromFamily(ctx context.Context, familyID int, userID int64) error {
	q := `DELETE FROM users_to_families WHERE user_id = $1 AND family_id = $2`

	_, err := ur.db.Exec(ctx, q, userID, familyID)
	if err != nil {
		ur.sl.Error("failed to leave family", slog.String("err", err.Error()))
		return err
	}

	return nil
}
