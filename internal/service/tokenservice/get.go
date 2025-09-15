package tokenservice

import (
	"context"
	"errors"
	"log/slog"
	"monofamily/internal/entity"

	"github.com/jackc/pgx/v4"
)

func (ts *TokenService) Get(ctx context.Context, familyID int, userID int64) (bool, *entity.UserBankToken, error) {

	hasToken := true

	ubt, err := ts.tokenProvider.Get(ctx, familyID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ts.sl.Debug("token not found")
			hasToken = false
		} else {
			ts.sl.Error("unable to get token from db", slog.String("err", err.Error()))
			return false, nil, err
		}
	}

	return hasToken, ubt, nil
}
