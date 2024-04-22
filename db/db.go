package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

var (
	// MigrationPath : DEFAULT DIRECTORY NAME
	MigrationPath = "schemas"

	// DriverName : DATABASE DRIVER NAME
	DriverName string

	// Host : DATABASE HOST NAME
	Host string

	// Port : DATABASE PORT NUMBER DEFAULT IS 5432
	Port = "5432"

	// User : DATABASE USER NAME
	User string

	// Password : DATABASE USER PASSWORD
	Password string

	// Name : DATABASE NAME WHERE DB OPERATIONS WILL BE PERFORMED
	Name string

	// DataSourceName : DSN GENERATED FROM ABOVE DB CREDENTIALS
	DataSourceName string
)

// InitMigrator returns migrate object after doing migrations to the database returned object can later be used to down the applied migrations
func InitMigrator() *migrate.Migrate {
	log.Println("Opening migration DB connection to " + DataSourceName)

	// OPENING DATABASE CONNECTION USING DSN AND DRIVER NAME
	newDB, openDBError := sql.Open(DriverName, DataSourceName)
	if openDBError != nil {
		log.Fatal("Error Opening Migration DB", openDBError)
	}
	log.Println("Fetching DB Driver")

	// CREATING DATABASE DRIVER USING NEWLY OPENED DATABASE CONNECTION
	DBDriver, dbDriverError := postgres.WithInstance(newDB, &postgres.Config{})
	if dbDriverError != nil {
		log.Fatal("Error Fetching DB Driver for Migration DB", dbDriverError)
	}
	log.Println("Initializing migrations")

	// CREATING NEW MIGRATOR USING DATABASE DRIVER INSTANCE, DATABASE NAME AND MIGRATION FILE PATH
	newMigrator, openMigratorError := migrate.NewWithDatabaseInstance(MigrationPath, Name, DBDriver)
	if openMigratorError != nil {
		log.Fatal("Error opening migrator: ", openMigratorError)
	}
	log.Println("Migrating to database")

	// MIGRATING ALL SCHEMA PRESENT IN MIGRATION PATH
	migrationError := newMigrator.Up()
	if migrationError != nil {
		if errors.Is(migrationError, migrate.ErrNoChange) {
			log.Println("Nothing to migrate, Skipping migration")
		} else {
			log.Fatal("Error migrating database: ", migrationError)
		}
	}
	log.Println("Migration completed successfully")
	log.Println("Closing Migration Database Connection")
	return newMigrator
}

// CleanupMigration drops all table present in migration and closes the migration database connection
func CleanupMigration(migratorObj *migrate.Migrate) {
	log.Println("Completed All Operations Dropping All Tables Created From Migration")

	// DROPPING ALL TABLES THAT WERE PREVIOUSLY CREATED AT THE START OF PROGRAM
	downMigrationError := migratorObj.Down()
	if downMigrationError != nil {
		log.Fatal("Error Dropping Migration: ", downMigrationError)
	}
	log.Println("Closing Migrator")

	// CLOSING MIGRATION AND DATABASE CONNECTION
	migrationCloseError, migrationDbCloseError := migratorObj.Close()
	if migrationDbCloseError != nil {
		log.Fatal("Error Closing Migration DB: ", migrationDbCloseError)
	}
	if migrationCloseError != nil {
		log.Fatal("Error Closing Migrator: ", migrationCloseError)
	}
	log.Println("Clean UP Successful. Program Exited.")
}
