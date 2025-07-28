package usecase

import (
	"context"
	"monofamily/internal/entity"
)

func (uc *UseCase) GetFamiliesByUserID(ctx context.Context, userID int64) ([]entity.Family, error) {
	return uc.familyService.GetFamiliesByUserID(ctx, userID)
}
