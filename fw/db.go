package fw

import "database/sql"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type DBConnector interface {
	Connect(config DBConfig) (*sql.DB, error)
}

type DBMigrationTool interface {
	MigrateUp(db *sql.DB, migrationRoot string) error
	MigrateDown(db *sql.DB, migrationRoot string) error
}
