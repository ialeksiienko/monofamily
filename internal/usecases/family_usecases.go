package usecases

import (
	"errors"
	"log/slog"
	"main-service/internal/entities"
	"main-service/internal/sessions"
	"main-service/internal/sl"
	"time"

	"github.com/jackc/pgx/v4"
)

type FamilyService struct {
	familyRepo FamilyRepository
	sl         *sl.MyLogger
}

func NewFamilyService(familyRepo FamilyRepository, sl *sl.MyLogger) *FamilyService {
	return &FamilyService{familyRepo: familyRepo, sl: sl}
}

func (s *FamilyService) Create(familyName string, userID int64) (string, time.Time, error) {
	f, err := s.familyRepo.CreateFamily(&entities.Family{
		Name:      familyName,
		CreatedBy: userID,
	})
	if err != nil {
		s.sl.Error("failed to create family", slog.Int("familyID", int(userID)), slog.String("err", err.Error()))
		return "", time.Time{}, err
	}

	s.sl.Debug("family created", slog.Int("familyID", f.ID))

	saveErr := s.familyRepo.SaveUserToFamily(f.ID, userID)
	if saveErr != nil {
		s.sl.Error("unable to save user to family", slog.Int("userID", int(userID)), slog.String("err", saveErr.Error()))
		return "", time.Time{}, err
	}

	code := generateInviteCode()

	expiresAt, err := s.familyRepo.SaveFamilyInviteCode(userID, f.ID, code)
	if err != nil {
		s.sl.Error("failed to save family invite code", slog.Int("familyID", int(userID)), slog.String("err", err.Error()))
		return "", time.Time{}, err
	}

	return code, expiresAt, nil
}

func (s *FamilyService) Join(code string, userID int64) (string, error) {
	f, expiresAt, err := s.familyRepo.GetFamilyByCode(code)
	if err != nil {
		s.sl.Error("failed to get family by code", slog.String("err", err.Error()))
		if errors.Is(err, pgx.ErrNoRows) {
			s.sl.Debug("family not found with code")
			return "", &CustomError[struct{}]{
				Msg:  "family not found by invite code",
				Code: ErrCodeFamilyNotFound,
			}
		}
		return "", err
	}

	if time.Now().After(expiresAt) {
		s.sl.Error("expired family code")
		return "", &CustomError[time.Time]{
			Data: expiresAt,
			Msg:  "family invite code expired",
			Code: ErrCodeFamilyCodeExpired,
		}
	}

	saveErr := s.familyRepo.SaveUserToFamily(f.ID, userID)
	if saveErr != nil {
		s.sl.Error("unable to save user to family", slog.Int("userID", int(userID)), slog.String("err", saveErr.Error()))
		return "", saveErr
	}

	return f.Name, nil
}

func (s *FamilyService) Enter(userID int64) ([]entities.Family, error) {
	families, err := s.familyRepo.GetFamiliesByUserID(userID)
	if err != nil {
		s.sl.Error("failed to get family by userID", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
		return nil, err
	}

	if len(families) == 0 {
		return nil, &CustomError[struct{}]{
			Msg:  "user has no family",
			Code: ErrCodeUserHasNoFamily,
		}
	}

	sessions.SetUserPageSession(userID, &sessions.UPSessions{
		Page:     0,
		Families: families,
	})

	return families, nil
}