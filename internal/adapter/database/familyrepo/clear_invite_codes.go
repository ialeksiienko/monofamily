package familyrepo

import (
	"context"
	"log/slog"
)

func (fr *FamilyRepository) ClearInviteCodes(ctx context.Context) error {
	_, err := fr.db.Exec(ctx, `
		DELETE FROM family_invite_codes
		WHERE expires_at < NOW()
	`)
	if err != nil {
		fr.sl.Error("failed to delete expired invite codes", slog.String("err", err.Error()))
		return err
	}

	return nil
}
