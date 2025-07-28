package usecase

import (
	"context"
	"errors"
	"monofamily/internal/entity"
	"monofamily/internal/errorsx"

	"github.com/jackc/pgx/v4"
)

func (uc *UseCase) SelectFamily(ctx context.Context, familyID int, userID int64) (bool, *entity.Family, error) {
	f, err := uc.familyService.GetFamilyByID(ctx, familyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, errorsx.NewError("family not found", errorsx.ErrCodeFamilyNotFound, struct{}{})
		}
		return false, nil, err
	}

	isAdmin := f.CreatedBy == userID

	return isAdmin, f, nil
}
