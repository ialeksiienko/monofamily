package usecases

import (
	"errors"
	"log/slog"
	"main-service/internal/entities"
	"strconv"

	"github.com/jackc/pgx/v4"
)

func (s *FamilyService) SelectFamily(userID int64, data string) (bool, *entities.Family, error) {
	familyID, err := strconv.Atoi(data)
	if err != nil {
		s.sl.Error("failed to convert family id to int", slog.String("error", err.Error()))
		return false, nil, err
	}

	f, err := s.familyRepo.GetFamilyBy("id", familyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, &CustomError[struct{}]{
				Msg: "family not found",
				Code: ErrCodeFamilyNotFound,
			}
		}
		return false, nil, err
	}

	isAdmin := f.CreatedBy == userID

	return isAdmin, f, nil
}