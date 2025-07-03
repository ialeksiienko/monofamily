package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (pool Database) SaveFamilyInviteCode(userID int64, familyID int, code string) (time.Time, error) {
	q := `INSERT INTO family_invite_codes 
    	(family_id, code, created_by, expires_at)
    	VALUES ($1, $2, $3, $4)
    	RETURNING expires_at`

	var expiresAt time.Time

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))
	err := pool.DB.QueryRow(context.Background(), q, familyID, code, userID, time.Now().UTC().Add(48*time.Hour)).Scan(&expiresAt)
	if err != nil {
		pool.logger.Error("failed to save family invite code", slog.String("err", err.Error()),
			slog.Int("user_id", int(userID)), slog.Int("family_id", familyID))
		return time.Time{}, err
	}

	return expiresAt, err
}

func (pool Database) ClearInviteCodes() error {
	ctx := context.Background()

	_, err := pool.DB.Exec(ctx, `
		DELETE FROM family_invite_codes
		WHERE expires_at < NOW()
	`)
	if err != nil {
		pool.logger.Error("failed to delete expired invite codes", slog.String("err", err.Error()))
		return err
	}

	return nil
}
