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

func (r *NodeRepository) GetAllNodes() ([]Node, error) {
	rows, err := r.db.Query("SELECT id, ip, ssh_port, ssh_password, status FROM nodes ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node
	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.ID, &n.IP, &n.SSHPort, &n.SSHPassword, &n.Status); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (r *NodeRepository) UpdateNode(id int, ip, port, password string) error {
	_, err := r.db.Exec("UPDATE nodes SET ip = ?, ssh_port = ?, ssh_password = ? WHERE id = ?", ip, port, password, id)
	return err
}

func (r *NodeRepository) DeleteNode(id int) error {
	_, err := r.db.Exec("DELETE FROM nodes WHERE id = ?", id)
	return err
}

func (r *NodeRepository) GetActiveNode() (string, string, error) {
	row := r.db.QueryRow("SELECT ip, ssh_port FROM nodes WHERE status = 'active' LIMIT 1")
	var ip, port string
	err := row.Scan(&ip, &port)
	return ip, port, err
}

func (r *NodeRepository) MarkNodeFailed(ip string) error {
	_, err := r.db.Exec("UPDATE nodes SET status = 'failed' WHERE ip = ?", ip)
	return err
}

func (r *NodeRepository) MarkNodeActive(id int) error {
	_, err := r.db.Exec("UPDATE nodes SET status = 'active' WHERE id = ?", id)
	return err
}
