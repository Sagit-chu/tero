// backend/repository/node_test.go
package repository

import (
	"database/sql"
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/db"
)

func TestNodeRepo(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewNodeRepository(database)

	_, err = repo.GetNextStandby()
	if err != sql.ErrNoRows {
		t.Fatalf("Expected ErrNoRows, got %v", err)
	}

	err = repo.AddNode("1.1.1.1", "22", "pass")
	if err != nil {
		t.Fatal(err)
	}

	node, err := repo.GetNextStandby()
	if err != nil {
		t.Fatal(err)
	}
	if node.IP != "1.1.1.1" {
		t.Fatalf("Expected IP 1.1.1.1, got %s", node.IP)
	}
}

func TestNodeRepo_GetAllNodes(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewNodeRepository(database)

	nodes, err := repo.GetAllNodes()
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 0 {
		t.Fatalf("Expected 0 nodes, got %d", len(nodes))
	}

	repo.AddNode("1.1.1.1", "22", "pass")
	repo.AddNode("2.2.2.2", "22", "pass")

	nodes, err = repo.GetAllNodes()
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(nodes))
	}
}
