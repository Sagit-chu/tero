// backend/repository/node.go
package repository

import "database/sql"

type Node struct {
	ID          int
	IP          string
	SSHPort     string
	SSHPassword string
	Status      string
}

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) AddNode(ip, port, password string) error {
	_, err := r.db.Exec("INSERT INTO nodes (ip, ssh_port, ssh_password) VALUES (?, ?, ?)", ip, port, password)
	return err
}

func (r *NodeRepository) GetNextStandby() (*Node, error) {
	row := r.db.QueryRow("SELECT id, ip, ssh_port, ssh_password, status FROM nodes WHERE status = 'standby' ORDER BY id ASC LIMIT 1")
	var n Node
	err := row.Scan(&n.ID, &n.IP, &n.SSHPort, &n.SSHPassword, &n.Status)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
