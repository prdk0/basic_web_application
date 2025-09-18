package dbrepo

import (
	"bookings/internals/config"
	"bookings/internals/repository"
	"database/sql"
)

type postgreDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgreDbRepo{
		App: a,
		DB:  conn,
	}
}
