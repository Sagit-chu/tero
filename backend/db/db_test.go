package db

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	os.Remove("test.db")
	db, err := InitDB("test.db")
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