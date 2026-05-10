// backend/repository/config_test.go
package repository

import (
	"database/sql"
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/db"
)

func TestConfigRepo(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewConfigRepository(database)

	err = repo.SetConfig("interval", "5m")
	if err != nil {
		t.Fatal(err)
	}

	val, err := repo.GetConfig("interval")
	if err != nil {
		t.Fatal(err)
	}
	if val != "5m" {
		t.Fatalf("Expected 5m, got %s", val)
	}

	_, err = repo.GetConfig("missing")
	if err != sql.ErrNoRows {
		t.Fatalf("Expected ErrNoRows, got %v", err)
	}
}
