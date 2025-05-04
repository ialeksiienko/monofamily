package repository

import (
	"context"
	"fmt"
	"log/slog"
	"monofinances/internal/models"
)

func (pool Database) CreateFamily(inp *models.Family) (*models.Family, error) {
	q := `INSERT INTO families (name) VALUES ($1) RETURNING id, name`

	f := new(models.Family)

	err := pool.DB.QueryRow(context.Background(), q, inp.Name).Scan(&f.ID, &f.Name)
	if err != nil {
		pool.logger.Error("failed to create family", slog.String("err", err.Error()), slog.String("family", inp.Name))
		return nil, err
	}

	return f, err
}

func (pool Database) GetFamilyByUserID(userID int64) ([]models.Family, error) {
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
