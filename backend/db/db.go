package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS config (key TEXT PRIMARY KEY, value TEXT);
	CREATE TABLE IF NOT EXISTS nodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		ssh_port TEXT NOT NULL,
		ssh_password TEXT NOT NULL,
		status TEXT DEFAULT 'standby',
		flvx_node_id INTEGER DEFAULT 0,
		flvx_node_name TEXT DEFAULT ''
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}