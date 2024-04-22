package conf

import (
	"embed"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sqlc-demo/db"
	"strings"
)

//go:embed *
var configFile embed.FS

// LoadEnvVariables loads the env variable from .env file and sets all env variable
func LoadEnvVariables() {

	// FETCHING ENV FILE FROM EMBED FS
	envFile, openEmbedFileError := configFile.Open(".env")
	if openEmbedFileError != nil {
		log.Fatal("Error Opening Embedded .env File: ", openEmbedFileError)
	}

	// PARSING ENV FILE FETCHED FROM EMBED FS
	envData, envParseError := godotenv.Parse(envFile)
	if envParseError != nil {
		log.Fatal("Error Parsing .env File: ", envParseError)
	}

	currentEnvInRaw := make(map[string]bool)
	rawEnv := os.Environ()
	// CHECKING IF KEY ALREADY EXIST IN ENV IF EXIST KEY IS MAPPED TRUE
	for _, rawEnvData := range rawEnv {
		key := strings.Split(rawEnvData, "=")[0]
		currentEnvInRaw[key] = true
	}

	// IF KEY VALUE FROM envData IS ALREADY IN RAW THEN IT DOES NOT SET ENV VALUE FETCHED FROM EMBED FS
	for key, value := range envData {
		if !currentEnvInRaw[key] {
			setenvError := os.Setenv(key, value)
			if setenvError != nil {
				log.Fatal("Error setting env vars: ", setenvError)
				return
			}
		}
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
