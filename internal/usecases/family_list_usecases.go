package usecases

import (
	"errors"
	"log/slog"
	"main-service/internal/sessions"
	"strconv"

	"github.com/jackc/pgx/v4"
)

func (s *FamilyService) SelectFamily(userID int64, data string) (bool, string, error) {
	familyID, err := strconv.Atoi(data)
	if err != nil {
		s.sl.Error("failed to convert family id to int", slog.String("error", err.Error()))
		return false, "", err
	}

	f, err := s.familyRepo.GetFamilyBy("id", familyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", &CustomError[struct{}]{
				Msg: "family not found",
				Code: ErrCodeFamilyNotFound,
			}
		}
		return false, "",err
	}

	isAdmin := f.CreatedBy == userID

	sessions.SetUserState(userID, &sessions.UserState{
		Family: f,
	})

	return isAdmin, f.Name, nil
}