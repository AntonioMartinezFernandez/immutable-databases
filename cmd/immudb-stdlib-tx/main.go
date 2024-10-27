package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

const (
	IMMUDB_USER            = "immudb"
	IMMUDB_PWD             = "immudb"
	IMMUDB_DB_NAME         = "defaultdb"
	IMMUDB_HOST            = "localhost"
	IMMUDB_PORT            = 3322
	IMMUDB_SSL             = false
	DB_TABLE               = "sql_posts"
	SQL_IMMUDB_DRIVER_NAME = "instrumented-sql-immudb"
)

var (
	maxLifetimeInMinutes = 10
	maxConnections       = 10
	connIdle             = 2
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleTxErr(err error, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		log.Fatal("Transaction rolled back: ", err)
	}
}

func main() {
	ctx := context.Background()
	startTime := time.Now()

	// Stablish database connection configuration
	opts := immudb.DefaultOptions().
		WithAddress(IMMUDB_HOST).
		WithPort(IMMUDB_PORT).
		WithUsername(IMMUDB_USER).
		WithPassword(IMMUDB_PWD).
		WithDatabase(IMMUDB_DB_NAME)

	// Create a new standard database/sql client and stablish pool connection
	dbClient := stdlib.OpenDB(opts)
	defer dbClient.Close() // Close the connection when the application is done

	dbClient.SetConnMaxLifetime(time.Duration(maxLifetimeInMinutes) * time.Minute)
	dbClient.SetMaxOpenConns(maxConnections)
	dbClient.SetMaxIdleConns(connIdle)

	pingErr := dbClient.Ping()
	handleErr(pingErr)

	// Create table if not exists
	dbCreateSqlQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER AUTO_INCREMENT,
			personId VARCHAR[128],
			text VARCHAR[4096],
			active BOOLEAN,
			PRIMARY KEY id
			);`,
		DB_TABLE,
	)
	sqlResult, err := dbClient.ExecContext(ctx, dbCreateSqlQuery)
	handleErr(err)
	fmt.Println("Table created with sqlResult: ", sqlResult)

	// Transactionally insert 1000 random rows into the table
	// !IMPORTANT: statements and sql results methods (like lastInsertId) are not supported by immudb
	tx, err := dbClient.BeginTx(ctx, nil)
	handleTxErr(err, tx)

	createPostSqlQuery := fmt.Sprintf(
		`INSERT INTO %s (personId, text, active) VALUES ($1, $2, $3)`,
		DB_TABLE,
	)

	for i := 0; i < 1000; i++ {
		_, err := tx.ExecContext(
			ctx,
			createPostSqlQuery,
			uuid.New().String(),
			faker.Paragraph(),
			i%2 == 0,
		)
		handleTxErr(err, tx)
	}

	txCommitErr := tx.Commit()
	handleTxErr(txCommitErr, tx)

	fmt.Println("SQL Transaction committed in: ", time.Since(startTime))
}
