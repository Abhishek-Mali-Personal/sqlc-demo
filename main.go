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
	LoadEnvVariables()
}

var (
	migrationPath  string
	driverName     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	dataSourceName string
)

func fetchAbsolutePath(relativePath string) (string, error) {

	// Get the absolute path
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}

	// Check if the file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return "", err
	}
	return absolutePath, nil
}

func LoadEnvVariables() {
	// Load .env file
	loadEnvError := godotenv.Load(".env")
	if loadEnvError != nil {
		log.Fatal("Error loading .env file", loadEnvError)
	}

	fetchAbsolutePath, fetchAbsolutePathError := fetchAbsolutePath("schemas")
	if fetchAbsolutePathError != nil {
		log.Fatal("Error Fetching Absolute Path: ", fetchAbsolutePathError)
	}
	migrationPath = "file:///" + fetchAbsolutePath
	driverName = os.Getenv("DRIVER_NAME")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
}

func InitMigrator() *migrate.Migrate {
	log.Println("Opening migration DB connection to " + dataSourceName)
	newDB, openDBError := sql.Open(driverName, dataSourceName)
	if openDBError != nil {
		log.Fatal("Error Opening Migration DB", openDBError)
	}
	log.Println("Fetching DB Driver")
	DBDriver, dbDriverError := postgres.WithInstance(newDB, &postgres.Config{})
	if dbDriverError != nil {
		log.Fatal("Error Fetching DB Driver for Migration DB", dbDriverError)
	}
	log.Println(migrationPath, strings.Index(migrationPath, ":"))
	log.Println("Initializing migrations")
	newMigrator, openMigratorError := migrate.NewWithDatabaseInstance(migrationPath, DBName, DBDriver)
	if openMigratorError != nil {
		log.Fatal("Error opening migrator: ", openMigratorError)
	}
	log.Println("Migrating to database")
	migrationError := newMigrator.Up()
	if migrationError != nil {
		log.Fatal("Error migrating database: ", migrationError)

	}
	log.Println("Migration completed successfully")
	log.Println("Closing Migration Database Connection")
	return newMigrator
}

func main() {
	migratorObj := InitMigrator()
	defer cleanupMigration(migratorObj)
	bgCtx := context.Background()
	log.Println("Creating database connection")
	conn, connError := pgx.Connect(bgCtx, dataSourceName)
	if connError != nil {
		log.Fatal("Error Connecting to database: ", connError)
	}
	log.Println("Connected to database")
	query := models.New(conn)
	log.Println("Executing Insert query")
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
	retrievedLookups, retrieveLookupError := query.ListLookupsByDisplayText(bgCtx, "GoLang")
	if retrieveLookupError != nil {
		log.Fatal("Error Retrieving Lookup: ", retrieveLookupError)
	}
	log.Println("Data Retrieved successfully")
	log.Printf("%#v", retrievedLookups)
	log.Println("Closing Database Connection")
	connCloseError := conn.Close(bgCtx)
	if connCloseError != nil {
		log.Fatal("Error Closing Database: ", connCloseError)
	}
	log.Println("Database connection closed")
	log.Println("Exiting Program")
}

func cleanupMigration(migratorObj *migrate.Migrate) {
	log.Println("Completed All Operations Dropping All Tables Created From Migration")
	downMigrationError := migratorObj.Down()
	if downMigrationError != nil {
		log.Fatal("Error Dropping Migration: ", downMigrationError)
	}
	log.Println("Closing Migrator")
	migrationCloseError, migrationDbCloseError := migratorObj.Close()
	if migrationDbCloseError != nil {
		log.Fatal("Error Closing Migration DB: ", migrationDbCloseError)
	}
	if migrationCloseError != nil {
		log.Fatal("Error Closing Migrator: ", migrationCloseError)
	}
	log.Println("Clean UP Successful. Program Exited.")
}
