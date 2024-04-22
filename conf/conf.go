package conf

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sqlc-demo/db"
	"sqlc-demo/service"
	"strings"
)

// LoadEnvVariables loads the env variable from .env file and sets all env variable
func LoadEnvVariables() {
	// LOAD .env FILE
	loadEnvError := godotenv.Load(".env")
	if loadEnvError != nil {
		log.Fatal("Error loading .env file", loadEnvError)
	}

	// FETCH ABSOLUTE PATH OF schemas DIRECTORY
	fetchAbsolutePath, fetchAbsolutePathError := service.FetchAbsolutePath(db.MigrationPath)
	if fetchAbsolutePathError != nil {
		log.Fatal("Error Fetching Absolute Path: ", fetchAbsolutePathError)
	}

	// ADD file:/// PREFIX IF NOT PRESENT IN ABSOLUTE PATH SPECIFIED AS MIGRATION FUNCTION NEEDS THIS FORMAT
	if !strings.HasPrefix(fetchAbsolutePath, "file:///") {
		db.MigrationPath = "file:///" + fetchAbsolutePath
	} else {
		db.MigrationPath = fetchAbsolutePath
	}

	// SET ALL ENV VARIABLES
	db.DriverName = os.Getenv("DRIVER_NAME")
	db.Host = os.Getenv("DB_HOST")
	db.Port = os.Getenv("DB_PORT")
	db.User = os.Getenv("DB_USER")
	db.Password = os.Getenv("DB_PASSWORD")
	db.Name = os.Getenv("DB_NAME")

	// SET DSN FOR DATABASE
	db.DataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)
}
