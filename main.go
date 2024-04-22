package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"sqlc-demo/conf"
	"sqlc-demo/db"
	"sqlc-demo/models"
	"time"
)

func init() {
	// LOAD ENV VARIABLES
	conf.LoadEnvVariables()
}

func main() {
	// INITIALISING DATABASE MIGRATION
	migratorObj := db.InitMigrator()

	// REVOKING ALL APPLIED MIGRATION AFTER COMPLETION OG PROGRAM
	defer db.CleanupMigration(migratorObj)

	// FETCHING BACKGROUND CONTEXT
	bgCtx := context.Background()
	log.Println("Creating database connection")

	// CONNECTING TO THE DATABASE FOR PERFORMING DATABASE OPERATIONS
	conn, connError := pgx.Connect(bgCtx, db.DataSourceName)
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
	log.Println("Executing Update query")

	// EXECUTING UPDATE QUERY USING CREATED QUERY OBJECT
	currentTime := time.Now()
	updatedLookup, updatedLookupError := query.UpdateLookup(bgCtx, models.UpdateLookupParams{
		IsActive:     false,
		UpdateDate:   &currentTime,
		UpdateUserID: pgtype.Int4{Int32: 1, Valid: true},
		LookupID:     createdLookup.LookupID,
	})
	if updatedLookupError != nil {
		log.Fatal("Error Updating Lookup: ", updatedLookupError)
	}
	log.Printf("Updated Lookup: %#v", updatedLookup)
	log.Println("Executing Retrieve query")

	// EXECUTING SELECT QUERY USING CREATED QUERY OBJECT
	retrievedLookups, retrieveLookupError := query.ListLookupsByDisplayText(bgCtx, updatedLookup.DisplayText)
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
