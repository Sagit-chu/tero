package db

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='nodes'").Scan(&name)
	if err != nil {
		t.Fatalf("Table nodes not created")
	}
}