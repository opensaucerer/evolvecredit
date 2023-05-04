package test

import (
	"log"
	"os"

	"github.com/opensaucerer/barf"
	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/global"
)

// Setup prepares the application for testing
// Setup could be made to use a separate environment file pointing to a test database
// such that the tables are dropped and recreated for each test
func Setup() {
	// load environment variables
	if err := barf.Env(global.ENV, os.Getenv("ENV_PATH")); err != nil {
		log.Fatal(err)
	}

	database.NewPostgreSQLConnection(global.ENV.PostgreSQLURI, global.ENV.PostgreSQLConnections)

	database.ReadFileAndExecuteQueries(global.ENV.SQLFilePath)
}

// Teardown cleans up the application after testing
// drop all tables in the database
func Teardown() {
	// database.DropAllTables()
}
