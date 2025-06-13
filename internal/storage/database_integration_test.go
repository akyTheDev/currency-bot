package storage

import (
	"fmt"
	"os"
	"testing"
)

func TestOpenDB_Success(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Errorf("DATABASE_URL is not set.")
	}

	db, err := OpenDB(dsn)
	if err != nil {
		t.Errorf("OpenDB failed: %v", err)
	}
	defer db.Close()

	// Double check ping.
	if err := db.Ping(); err != nil {
		t.Errorf("Ping to real PSQL failed: %v", err)
	}
}

func TestOpenDB_Fail(t *testing.T) {
	dsn := "postgres://username:password@localhost:5432/wrongDB?sslmode=disable"

	_, err := OpenDB(dsn)
	if err == nil {
		t.Fatalf("OpenDB didn't return the expected error: %v", err)
	}
	fmt.Print(err)
}
