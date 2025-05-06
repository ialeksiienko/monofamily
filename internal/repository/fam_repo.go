package repository

import (
	"context"
	"fmt"
	"log/slog"
	"main-service/internal/models"
	"time"
)

func (pool Database) CreateFamily(inp *models.Family) (*models.Family, error) {
	q := `INSERT INTO families (user_id, name) 
			VALUES ($1, $2) RETURNING id, user_id, name`

	f := new(models.Family)

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))
	err := pool.DB.QueryRow(context.Background(), q, inp.CreatedBy, inp.Name).Scan(&f.ID, &f.CreatedBy, &f.Name)
	if err != nil {
		pool.logger.Error("failed to create family", slog.String("err", err.Error()), slog.String("family", inp.Name))
		return nil, err
	}

	return f, err
}

func (pool Database) GetFamiliesByUserID(userID int64) ([]models.Family, error) {
	q := `SELECT f.id, f.name 
	FROM users_to_families utf
	JOIN families f ON f.id = utf.family_id
	WHERE utf.user_id = $1`

	rows, err := pool.DB.Query(context.Background(), q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))

	var families []models.Family
	for rows.Next() {
		var f models.Family
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			return nil, err
		}
		families = append(families, f)
	}

	if err := rows.Err(); err != nil {
		pool.logger.Error("failed to get family by user id", slog.String("err", err.Error()))
		return nil, err
	}

	return families, nil
}

func (pool Database) SaveFamilyInviteCode(userId int64, familyId int, code string) (time.Time, error) {
	q := `INSERT INTO family_invites 
    	(family_id, code, created_by, expires_at)
    	VALUES ($1, $2, $3, $4)
    	RETURNING expires_at`

	var expiresAt time.Time

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))
	err := pool.DB.QueryRow(context.Background(), q, familyId, code, userId, time.Now().Add(48*time.Hour)).Scan(&expiresAt)
	if err != nil {
		pool.logger.Error("failed to save family invite code", slog.String("err", err.Error()),
			slog.Int("user_id", int(userId)), slog.Int("family_id", familyId))
		return time.Time{}, err
	}

	return expiresAt, err
}

func (pool Database) GetFamilyByCode(code string) (*models.Family, time.Time, error) {
	q := `SELECT f.id, f.created_by, f.name, fi.expires_at
		FROM family_invites fi
		JOIN families f ON f.id = fi.family_id
		WHERE fi.code = $1`

	f := new(models.Family)

	var expiresAt time.Time

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))
	err := pool.DB.QueryRow(context.Background(), q, code).Scan(&f.ID, &f.CreatedBy, &f.Name, &expiresAt)
	if err != nil {
		pool.logger.Error("failed to get family", slog.String("err", err.Error()))
		return nil, time.Time{}, err
	}

	return f, expiresAt, nil
}

func (pool Database) SaveUserToFamily(f *models.Family) error {
	q := `INSERT INTO users_to_families (user_id, family_id)
			VALUES ($1, $2)`

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))
	_, err := pool.DB.Exec(context.Background(), q, f.CreatedBy, f.ID)
	if err != nil {
		pool.logger.Error("failed to create family",
			slog.String("err", err.Error()),
			slog.Int("family_id", f.ID))
		return err
	}

	return nil
}
