package familyservice

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

type WithTransaction interface {
	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
}

func (s *FamilyService) DeleteFamily(ctx context.Context, familyID int) error {
	err := s.withTransaction.WithTransaction(ctx, func(tx pgx.Tx) error {
		return s.familyDeletor.DeleteFamily(ctx, tx, familyID)
	})
	if err != nil {
		s.sl.Error("failed to delete family", slog.Int("family_id", familyID), slog.String("error", err.Error()))
		return err
	}

	return nil
}
