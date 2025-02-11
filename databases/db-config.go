package databases

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgresRepo struct{}

// NewPostgresRepo func
func NewPostgresRepo() DatabaseRepo {
	return &postgresRepo{}
}

// Connect func
func (*postgresRepo) Connect(host string, port int, user string, password string, dbName string, dbMigrateVersion uint, runMigration bool, dbDriver string) (*sqlx.DB, error) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbName,
		password,
	)

	dbConn, err := sqlx.Open(dbDriver, connectionStr)
	if err != nil {
		log.Println("Error Postgres Database Connection..")
		log.Println(err)
		os.Exit(3)
	}

	dbConn.SetConnMaxLifetime(time.Minute * 3)
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(2)

	err = dbConn.Ping()
	if err != nil {
		log.Println("Could't Connect to Postgres Database..")
		os.Exit(3)
	}

	log.Println("Postgres Database Connected..")

	if runMigration {
		if err = runMigrationDB(dbConn, dbMigrateVersion, dbName, dbDriver); err != nil {
			log.Println("Failed to Migrate Postgres Database..")
			os.Exit(3)
		}
	}

	return dbConn, nil
}

func runMigrationDB(db *sqlx.DB, desiredVersion uint, dbName string, dbDriver string) error {
	log.Println("Postgres Database Migrating..")

	connection := sqlx.NewDb(db.DB, dbDriver)

	instance, err := postgres.WithInstance(
		connection.DB,
		&postgres.Config{
			MigrationsTable:       postgres.DefaultMigrationsTable,
			MultiStatementEnabled: true,
			DatabaseName:          dbName,
			SchemaName:            "public",
		},
	)
	if err != nil {
		log.Println("[err] [postgresRepo] [Migrate] [WithInstance] => ", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./databases/migrations",
		dbDriver, instance)
	if err != nil {
		log.Println("[err] [postgresRepo] [Migrate] [NewWithDatabaseInstance] => ", err)
		return err
	}

	if err = m.Force(-1); err != nil {
		log.Println("[err] [postgresRepo] [Migrate] [Force] => ", err)
		return err
	}

	if err = m.Migrate(desiredVersion); err != nil {
		log.Println("[err] [postgresRepo] [Migrate] [Migrate] => ", err)
		return err
	}

	log.Println("Postgres Database Migrated..")

	return nil
}
