package repository

import (
	"context"
	"fmt"
	"log/slog"
	"main-service/internal/models"
)

func (pool Database) SaveUser(user *models.User) (*models.User, error) {
	q := `INSERT INTO users (id, username, firstname)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE 
    		SET id = EXCLUDED.id
			RETURNING id, username, firstname, joined_at;`

	u := new(models.User)

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))
	err := pool.DB.QueryRow(context.Background(), q, user.ID, user.Username, user.Firstname).Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt)
	if err != nil {
		pool.logger.Error("failed to save user", slog.String("error", err.Error()))
		return nil, err
	}

	return u, nil
}

func (pool Database) GetAllUsersInFamily(f *models.Family) ([]models.User, error) {
	q := `SELECT u.username, u.firstname
	FROM users_to_families utf
	JOIN users u ON u.id = utf.user_id
	WHERE utf.family_id = $1`

	rows, err := pool.DB.Query(context.Background(), q, f.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pool.logger.Debug(fmt.Sprintf("SQL: %s", q))

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Username, &u.Firstname); err != nil {
			return nil, err
		}
		users = append(users, f)
	}

	if err := rows.Err(); err != nil {
		pool.logger.Error("failed to get family by user id", slog.String("err", err.Error()))
		return nil, err
	}

	return users, nil
}
