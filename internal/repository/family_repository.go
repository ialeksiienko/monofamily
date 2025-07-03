package repository

import (
	"context"
	"log/slog"
	"main-service/internal/entities"
	"time"
)

func (pool Database) CreateFamily(inp *entities.Family) (*entities.Family, error) {
	q := `INSERT INTO families (created_by, name) 
			VALUES ($1, $2) RETURNING id, created_by, name`

	f := new(entities.Family)

	err := pool.DB.QueryRow(context.Background(), q, inp.CreatedBy, inp.Name).Scan(&f.ID, &f.CreatedBy, &f.Name)
	if err != nil {
		pool.logger.Error("failed to create family", slog.String("err", err.Error()), slog.String("family", inp.Name))
		return nil, err
	}

	return f, err
}

func (pool Database) GetFamiliesByUserID(userID int64) ([]entities.Family, error) {
	q := `SELECT f.id, f.name 
	FROM users_to_families utf
	JOIN families f ON f.id = utf.family_id
	WHERE utf.user_id = $1`

	rows, err := pool.DB.Query(context.Background(), q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var families []entities.Family
	for rows.Next() {
		var f entities.Family
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

func (pool Database) GetFamilyByCode(code string) (*entities.Family, time.Time, error) {
	q := `SELECT f.id, f.created_by, f.name, fi.expires_at
		FROM family_invite_codes fi
		JOIN families f ON f.id = fi.family_id
		WHERE fi.code = $1`

	f := new(entities.Family)

	var expiresAt time.Time

	err := pool.DB.QueryRow(context.Background(), q, code).Scan(&f.ID, &f.CreatedBy, &f.Name, &expiresAt)
	if err != nil {
		pool.logger.Error("failed to get family", slog.String("err", err.Error()))
		return nil, time.Time{}, err
	}

	return f, expiresAt, nil
}

func (pool Database) GetFamilyByID(id int) (*entities.Family, error) {
	q := `SELECT id, created_by, name FROM families WHERE id = $1`

	f := new(entities.Family)
	err := pool.DB.QueryRow(context.Background(), q, id).Scan(&f.ID, &f.CreatedBy, &f.Name)
	if err != nil {
		pool.logger.Error("failed to get family", slog.String("err", err.Error()))
		return nil, err
	}

	return f, nil
}

func (pool Database) DeleteFamily(familyID int) error {
	ctx := context.Background()

	tx, err := pool.DB.Begin(ctx)
	if err != nil {
		pool.logger.Error("failed to begin transaction", slog.String("err", err.Error()))
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM users_to_families WHERE family_id = $1`, familyID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM family_invite_codes WHERE family_id = $1`, familyID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM families WHERE id = $1`, familyID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		pool.logger.Error("failed to commit transaction", slog.String("err", err.Error()))
		return err
	}

	return nil
}
