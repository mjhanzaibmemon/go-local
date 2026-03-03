package db

import "fmt"

type (
	// ConnectionArgs holds the parameters for establishing a DB connection (without DB name).
	ConnectionArgs struct {
		Username string
		Password string
		Endpoint string
		Port     string
	}

	// DSN represents a Data Source Name builder with extra flags for TLS and multi-statements.
	DSN struct {
		ConnectionArgs
		DatabaseName string
		// AllowMultiStatements allows running multiple statements in single query. Useful for migrations.
		AllowMultiStatements bool
		// TLSPreferred enables preferred TLS mode for the MySQL driver.
		TLSPreferred bool
	}
)

// NewDSN constructs a DSN instance for the provided connection arguments and DB name.
func NewDSN(connection *ConnectionArgs, databaseName string) *DSN {
	return &DSN{
		ConnectionArgs:       *connection,
		DatabaseName:         databaseName,
		AllowMultiStatements: true,
		TLSPreferred:         true,
	}
}

// String builds the DSN string with driver flags based on the configuration.
func (d *DSN) String() string {
	dsn := fmt.Sprintf(
		"%s%s?parseTime=true&loc=UTC&charset=utf8mb4",
		d.ConnectionArgs.String(),
		d.DatabaseName,
	)

	if d.TLSPreferred {
		dsn += "&tls=preferred"
	}

	if d.AllowMultiStatements {
		dsn += "&multiStatements=true"
	}

	return dsn
}

// String renders the connection arguments into the DSN prefix without the database name.
func (d *ConnectionArgs) String() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/",
		d.Username,
		d.Password,
		d.Endpoint,
		d.Port,
	)

	return dsn
}
