package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/luna-duclos/instrumentedsql"

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

func main() {
	ctx := context.Background()

	// even though the server address and port are defaults, setting them as a reference
	opts := immudb.DefaultOptions().
		WithAddress("localhost").
		WithPort(3322)

	client := immudb.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := client.OpenSession(
		ctx,
		[]byte(IMMUDB_USER),
		[]byte(IMMUDB_PWD),
		IMMUDB_DB_NAME)
	handleErr(err)

	// ensure connection is closed
	defer client.CloseSession(ctx)

	// Bonus: Here's how you instrument SQLX
	logger := instrumentedsql.LoggerFunc(
		func(ctx context.Context, msg string, keyvals ...interface{}) {
			log.Printf("%s %v", msg, keyvals)
		},
	)

	// Register the new driver as a wrapper of the standard immudb driver
	sql.Register(
		SQL_IMMUDB_DRIVER_NAME,
		instrumentedsql.WrapDriver(
			&stdlib.Driver{},
			instrumentedsql.WithLogger(logger),
		),
	)

	// Create a new standard database/sql client
	dbClient, err := sql.Open(
		SQL_IMMUDB_DRIVER_NAME,
		fmt.Sprintf("immudb://%s:%s@%s:%d/%s?sslmode=disable", IMMUDB_USER, IMMUDB_PWD, IMMUDB_HOST, IMMUDB_PORT, IMMUDB_DB_NAME),
	)
	handleErr(err)

	pingErr := dbClient.Ping()
	handleErr(pingErr)

	dbClient.SetConnMaxLifetime(time.Duration(maxLifetimeInMinutes) * time.Minute)
	dbClient.SetMaxOpenConns(maxConnections)
	dbClient.SetMaxIdleConns(connIdle)

	// Create our test table
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

	// Insert 1000 random values rows into the table
	for i := 0; i < 1000; i++ {
		sqlResult, err := dbClient.ExecContext(
			ctx,
			fmt.Sprintf(
				`INSERT INTO %s (personId, text, active) VALUES ($1, $2, $3)`, DB_TABLE),
			uuid.New().String(),
			faker.Paragraph(),
			i%2 == 0,
		)
		handleErr(err)
		rowsAffected, _ := sqlResult.RowsAffected()
		lastInsertId, _ := sqlResult.LastInsertId()
		fmt.Println("Inserted row with rowsAffected: ", rowsAffected, " and lastInsertId: ", lastInsertId)
	}

	// Close the connection
	dbClient.Close()
}
