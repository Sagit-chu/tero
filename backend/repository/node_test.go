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

func TestNodeRepo_UpdateAndDelete(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewNodeRepository(database)

	repo.AddNode("1.1.1.1", "22", "pass")
	nodes, _ := repo.GetAllNodes()
	id := nodes[0].ID

	// Test Update
	err = repo.UpdateNode(id, "2.2.2.2", "2222", "newpass")
	if err != nil {
		t.Fatal(err)
	}
	nodes, _ = repo.GetAllNodes()
	if nodes[0].IP != "2.2.2.2" {
		t.Fatalf("Expected IP 2.2.2.2, got %s", nodes[0].IP)
	}

	// Test Delete
	err = repo.DeleteNode(id)
	if err != nil {
		t.Fatal(err)
	}
	nodes, _ = repo.GetAllNodes()
	if len(nodes) != 0 {
		t.Fatalf("Expected 0 nodes, got %d", len(nodes))
	}
}
