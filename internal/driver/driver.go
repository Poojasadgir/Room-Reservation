package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL connects to a SQL database using the provided DSN and returns a pointer to a DB object and an error.
// The DB object contains a connection pool with the maximum number of open connections, maximum number of idle connections, and maximum lifetime of a connection set.
// If the connection to the database fails, it returns an error.
func ConnectSQL(dsn string) (*DB, error) {
	dbPool, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	dbPool.SetMaxOpenConns(maxOpenDbConn)
	dbPool.SetMaxIdleConns(maxIdleDbConn)
	dbPool.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = dbPool

	err = testDB(dbPool)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB tests the connection to the database by pinging it.
// It takes a pointer to a sql.DB object as input and returns an error if the ping fails.
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database connection using the provided DSN and returns a pointer to the sql.DB object.
// It returns an error if there is an issue with opening or pinging the database.
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
