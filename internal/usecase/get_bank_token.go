package usecase

import "context"

func (uc *UseCase) GetBankToken(ctx context.Context, familyID int, userID int64, token string) (string, error) {
	ubt, err := uc.tokenService.Save(ctx, familyID, userID, token)
	if err != nil {
		return "", err
	}

	decryptedToken, err := uc.tokenService.Decrypt(ubt.Token)
	if err != nil {
		return "", err
	}

	return decryptedToken, nil
}
