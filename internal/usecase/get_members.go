package usecase

import (
	"context"
	"monofamily/internal/entity"
	"monofamily/internal/service/userservice"
)

func (uc *UseCase) GetFamilyMembersInfo(ctx context.Context, family *entity.Family, userID int64) ([]userservice.MemberInfo, error) {
	return uc.userService.GetFamilyMembersInfo(ctx, family, userID)
}
