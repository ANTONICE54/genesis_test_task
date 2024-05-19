package db

import "database/sql"

// PostgresDB provides functions to execute SQL queries
type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB(db *sql.DB) PostgresDB {
	return PostgresDB{
		DB: db,
	}
}
