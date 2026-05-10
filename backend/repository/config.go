// backend/repository/config.go
package repository

import "database/sql"

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) SetConfig(key, value string) error {
	_, err := r.db.Exec("INSERT INTO config (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value", key, value)
	return err
}

func (r *ConfigRepository) GetConfig(key string) (string, error) {
	row := r.db.QueryRow("SELECT value FROM config WHERE key = ?", key)
	var val string
	err := row.Scan(&val)
	return val, err
}
