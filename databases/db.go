package databases

import "github.com/jmoiron/sqlx"

// DatabaseRepo interface
type DatabaseRepo interface {
	Connect(host string, port int, user string, password string, dbName string, dbMigrateVersion uint, runMigration bool, dbDriver string) (*sqlx.DB, error)
}
