package usecase

import (
	"context"
	"monofamily/internal/entity"
)

func (uc *UseCase) RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return uc.userService.Register(ctx, user)
}
