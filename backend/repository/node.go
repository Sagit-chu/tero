package repository

import "database/sql"

type Node struct {
	ID           int    `json:"ID"`
	IP           string `json:"IP"`
	SSHPort      string `json:"SSHPort"`
	SSHPassword  string `json:"SSHPassword"`
	Status       string `json:"Status"`
	FlvxNodeID   int    `json:"FlvxNodeID"`
	FlvxNodeName string `json:"FlvxNodeName"`
}

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) AddNode(ip, port, password string, flvxNodeID int, flvxNodeName string) error {
	_, err := r.db.Exec("INSERT INTO nodes (ip, ssh_port, ssh_password, flvx_node_id, flvx_node_name) VALUES (?, ?, ?, ?, ?)", ip, port, password, flvxNodeID, flvxNodeName)
	return err
}

func (r *NodeRepository) GetNextStandby() (*Node, error) {
	row := r.db.QueryRow("SELECT id, ip, ssh_port, ssh_password, status, COALESCE(flvx_node_id, 0), COALESCE(flvx_node_name, '') FROM nodes WHERE status = 'standby' ORDER BY id ASC LIMIT 1")
	var n Node
	err := row.Scan(&n.ID, &n.IP, &n.SSHPort, &n.SSHPassword, &n.Status, &n.FlvxNodeID, &n.FlvxNodeName)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NodeRepository) GetAllNodes() ([]Node, error) {
	rows, err := r.db.Query("SELECT id, ip, ssh_port, ssh_password, status, COALESCE(flvx_node_id, 0), COALESCE(flvx_node_name, '') FROM nodes ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node
	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.ID, &n.IP, &n.SSHPort, &n.SSHPassword, &n.Status, &n.FlvxNodeID, &n.FlvxNodeName); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (r *NodeRepository) UpdateNode(id int, ip, port, password string, flvxNodeID int, flvxNodeName string) error {
	_, err := r.db.Exec("UPDATE nodes SET ip = ?, ssh_port = ?, ssh_password = ?, flvx_node_id = ?, flvx_node_name = ? WHERE id = ?", ip, port, password, flvxNodeID, flvxNodeName, id)
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