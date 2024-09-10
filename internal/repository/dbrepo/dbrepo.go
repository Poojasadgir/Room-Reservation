package dbrepo

import (
	"database/sql"

	"github.com/Poojasadgir/room-reservation/internal/config"
	"github.com/Poojasadgir/room-reservation/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates a new instance of the postgresDBRepo struct that implements the DatabaseRepo interface.
// It takes a pointer to a sql.DB and a pointer to a config.AppConfig as arguments and returns a pointer to the DatabaseRepo interface.
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

// NewTestingsRepo creates a new instance of testDBRepo and returns it as a DatabaseRepo interface.
// It takes an instance of AppConfig as a parameter.
func NewTestingsRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
