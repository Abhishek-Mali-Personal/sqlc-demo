package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"sqlc-demo/models"
	"strings"
	"time"
)

func init() {
	// LOAD ENV VARIABLES
	LoadEnvVariables()
}

var (
	// DEFAULT DIRECTORY NAME
	migrationPath = "schemas"

	// DATABASE DRIVER NAME
	driverName string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// DSN GENERATED FROM ABOVE DB CREDENTIALS
	dataSourceName string
)

// fetchAbsolutePath returns an absolute path if file or directory exists else returns error
func fetchAbsolutePath(relativePath string) (string, error) {
	// GET THE ABSOLUTE PATH
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}

	// CHECK IF THE FILE EXISTS
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return "", err
	}
	return absolutePath, nil
}

// LoadEnvVariables loads the env variable from .env file and sets all env variable
func LoadEnvVariables() {
	// LOAD .env FILE
	loadEnvError := godotenv.Load(".env")
	if loadEnvError != nil {
		log.Fatal("Error loading .env file", loadEnvError)
	}

	// FETCH ABSOLUTE PATH OF schemas DIRECTORY
	fetchAbsolutePath, fetchAbsolutePathError := fetchAbsolutePath(migrationPath)
	if fetchAbsolutePathError != nil {
		log.Fatal("Error Fetching Absolute Path: ", fetchAbsolutePathError)
	}

	// ADD file:/// PREFIX IF NOT PRESENT IN ABSOLUTE PATH SPECIFIED AS MIGRATION FUNCTION NEEDS THIS FORMAT
	if !strings.HasPrefix(fetchAbsolutePath, "file:///") {
		migrationPath = "file:///" + fetchAbsolutePath
	} else {
		migrationPath = fetchAbsolutePath
	}

	// SET ALL ENV VARIABLES
	driverName = os.Getenv("DRIVER_NAME")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")

	// SET DSN FOR DATABASE
	dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
}

// InitMigrator returns migrate object after doing migrations to the database returned object can later be used to down the applied migrations
func InitMigrator() *migrate.Migrate {
	log.Println("Opening migration DB connection to " + dataSourceName)

	// OPENING DATABASE CONNECTION USING DSN AND DRIVER NAME
	newDB, openDBError := sql.Open(driverName, dataSourceName)
	if openDBError != nil {
		log.Fatal("Error Opening Migration DB", openDBError)
	}
	log.Println("Fetching DB Driver")

	// CREATING DATABASE DRIVER USING NEWLY OPENED DATABASE CONNECTION
	DBDriver, dbDriverError := postgres.WithInstance(newDB, &postgres.Config{})
	if dbDriverError != nil {
		log.Fatal("Error Fetching DB Driver for Migration DB", dbDriverError)
	}
	log.Println(migrationPath, strings.Index(migrationPath, ":"))
	log.Println("Initializing migrations")

	// CREATING NEW MIGRATOR USING DATABASE DRIVER INSTANCE, DATABASE NAME AND MIGRATION FILE PATH
	newMigrator, openMigratorError := migrate.NewWithDatabaseInstance(migrationPath, DBName, DBDriver)
	if openMigratorError != nil {
		log.Fatal("Error opening migrator: ", openMigratorError)
	}
	log.Println("Migrating to database")

	// MIGRATING ALL SCHEMA PRESENT IN MIGRATION PATH
	migrationError := newMigrator.Up()
	if migrationError != nil {
		log.Fatal("Error migrating database: ", migrationError)

	}
	log.Println("Migration completed successfully")
	log.Println("Closing Migration Database Connection")
	return newMigrator
}

func main() {
	// INITIALISING DATABASE MIGRATION
	migratorObj := InitMigrator()

	// REVOKING ALL APPLIED MIGRATION AFTER COMPLETION OG PROGRAM
	defer cleanupMigration(migratorObj)

	// FETCHING BACKGROUND CONTEXT
	bgCtx := context.Background()
	log.Println("Creating database connection")

	// CONNECTING TO THE DATABASE FOR PERFORMING DATABASE OPERATIONS
	conn, connError := pgx.Connect(bgCtx, dataSourceName)
	if connError != nil {
		log.Fatal("Error Connecting to database: ", connError)
	}
	log.Println("Connected to database")

	// CREATING QUERY OBJECT FROM CONNECTED DATABASE
	query := models.New(conn)
	log.Println("Executing Insert query")

	// EXECUTING INSERT QUERY USING CREATED QUERY OBJECT
	createdLookup, createLookupError := query.NewLookupWithConcurrencyKey(bgCtx, models.NewLookupWithConcurrencyKeyParams{
		TableName:      "Department",
		DisplayOrder:   0,
		DisplayText:    "GoLang",
		IsActive:       true,
		InternalKey:    "test",
		ConcurrencyKey: "test",
		CreateDate:     time.Now(),
		CreateUserID:   1,
		ValueText:      "Golang",
	})
	if createLookupError != nil {
		log.Fatal("Error Creating Lookup: ", createLookupError)
	}
	log.Printf("%#v", createdLookup)
	log.Println("Executing Retrieve query")

	// EXECUTING SELECT QUERY USING CREATED QUERY OBJECT
	retrievedLookups, retrieveLookupError := query.ListLookupsByDisplayText(bgCtx, "GoLang")
	if retrieveLookupError != nil {
		log.Fatal("Error Retrieving Lookup: ", retrieveLookupError)
	}
	log.Println("Data Retrieved successfully")
	log.Printf("%#v", retrievedLookups)
	log.Println("Closing Database Connection")

	// CLOSING DATABASE CONNECTION
	connCloseError := conn.Close(bgCtx)
	if connCloseError != nil {
		log.Fatal("Error Closing Database: ", connCloseError)
	}
	log.Println("Database connection closed")
	log.Println("Exiting Program")
}

// cleanupMigration drops all table present in migration and closes the migration database connection
func cleanupMigration(migratorObj *migrate.Migrate) {
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
