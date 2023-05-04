package types

type Env struct {
	// Port for the server to listen on
	Port string `barfenv:"key=PORT;required=true"`
	// Database connection string
	PostgreSQLURI string `barfenv:"key=POSTGRESQL_URI;required=true"`
	// Number of connections to the database
	PostgreSQLConnections int32 `barfenv:"key=POSTGRESQL_CONNECTIONS;required=true"`
	// Application token
	AppToken string `barfenv:"key=APP_TOKEN;required=true"`
	// Path to SQL file containing queries to be executed on startup
	SQLFilePath string `barfenv:"key=SQL_FILE_PATH;required=true"`
}
