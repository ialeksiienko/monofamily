package userrepo

import (
	"monofamily/internal/adapter/database"
	"monofamily/internal/pkg/sl"
)

type databaseInface interface {
	database.PgxIface
}

type UserRepository struct {
	db databaseInface
	sl sl.Logger
}

func New(db databaseInface, sl sl.Logger) *UserRepository {
	return &UserRepository{
		db: db,
		sl: sl,
	}
}
