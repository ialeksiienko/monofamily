package usecases

import (
	"errors"
	"main-service/internal/entities"

	"github.com/jackc/pgx/v4"
)

func (s *FamilyService) SelectFamily(familyID int, userID int64) (bool, *entities.Family, error) {
	f, err := s.familyProvider.GetFamilyByID(familyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, &CustomError[struct{}]{
				Msg:  "family not found",
				Code: ErrCodeFamilyNotFound,
			}
		}
		return false, nil, err
	}

	isAdmin := f.CreatedBy == userID

	return isAdmin, f, nil
}
