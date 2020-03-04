package testutil

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// DatabaseConnectionConfig holds configuration parameters for the database
// connection.
type DatabaseConnectionConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	Schema   string
}

// NewDatabaseConnectionConfig returns a new config instance holding default
// values.
func NewDatabaseConnectionConfig() *DatabaseConnectionConfig {
	return &DatabaseConnectionConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		User:     "postgres",
		Password: "password",
		Schema:   "theme_park",
	}
}

// NewDatabaseConnection returns a new *sqlx.DB instance connected
// to a local postgres database.
func NewDatabaseConnection(config *DatabaseConnectionConfig) (*sqlx.DB, error) {

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable search_path=%s",
		config.Host, config.Port, config.Database, config.User, config.Password, config.Schema,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewDatabaseConnectionFromEnv is similar to NewDatabaseConnection but grabs
// the required values from the environment.
func NewDatabaseConnectionFromEnv() (*sqlx.DB, error) {
	return nil, nil
}
