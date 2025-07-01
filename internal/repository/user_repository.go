package repository

import (
	"context"
	"log/slog"
	"main-service/internal/entities"
)

func (pool Database) SaveUser(user *entities.User) (*entities.User, error) {
	q := `INSERT INTO users (id, username, firstname)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE 
    		SET id = EXCLUDED.id
			RETURNING id, username, firstname, joined_at;`

	u := new(entities.User)

	err := pool.DB.QueryRow(context.Background(), q, user.ID, user.Username, user.Firstname).Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt)
	if err != nil {
		pool.logger.Error("failed to save user", slog.String("error", err.Error()))
		return nil, err
	}

	return u, nil
}

func (pool Database) GetAllUsersInFamily(familyID int) ([]entities.User, error) {
	q := `SELECT u.id, u.username, u.firstname, u.joined_at
	FROM users_to_families utf
	JOIN users u ON u.id = utf.user_id
	WHERE utf.family_id = $1`

	rows, err := pool.DB.Query(context.Background(), q, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var u entities.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Firstname, &u.JoinedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		pool.logger.Error("failed to get family by user id", slog.String("err", err.Error()))
		return nil, err
	}

	return users, nil
}

func (pool Database) SaveUserToFamily(familyID int, userID int64) error {
	q := `INSERT INTO users_to_families (user_id, family_id)
			VALUES ($1, $2)`

	_, err := pool.DB.Exec(context.Background(), q, userID, familyID)
	if err != nil {
		pool.logger.Error("failed to create family",
			slog.String("err", err.Error()),
			slog.Int("family_id", familyID))
		return err
	}

	return nil
}

func (pool Database) DeleteUserFromFamily(familyID int, userID int64) error {
	q := `DELETE FROM users_to_families WHERE user_id = $1 AND family_id = $2`

	_, err := pool.DB.Exec(context.Background(), q, userID, familyID)
	if err != nil {
		pool.logger.Error("failed to leave family", slog.String("err", err.Error()))
		return err
	}

	return nil
}
