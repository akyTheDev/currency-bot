package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var sqlOpen = sql.Open

func OpenDB(databaseURL string) (*sql.DB, error) {
	db, err := sqlOpen("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to open PSQL connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping PSQL: %w", err)
	}

	return db, nil
}
