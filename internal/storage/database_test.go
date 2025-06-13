package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func openDBWithMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	mock.ExpectPing().WillReturnError(nil)
	return db, mock
}

func TestOpenDB(t *testing.T) {
	origOpen := sqlOpen
	defer func() { sqlOpen = origOpen }()

	t.Run("Success", func(t *testing.T) {
		dbMock, mock := openDBWithMock(t)

		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			if driverName != "pgx" {
				return nil, fmt.Errorf("unexpected driver: %s", driverName)
			}
			return dbMock, nil
		}

		db, err := OpenDB("irrelevant-for-mock")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		defer db.Close()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		openDBWithMock(t)

		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New("ERROR")
		}

		_, err := OpenDB("irrelevant-for-mock")
		if err == nil {
			t.Fatal("Expected OpenDB to return an error when OpenDB function fails, but got nil.")
		}
	})
}
