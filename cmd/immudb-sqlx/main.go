package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
)

const (
	IMMUDB_USER            = "immudb"
	IMMUDB_PWD             = "immudb"
	IMMUDB_DB_NAME         = "defaultdb"
	IMMUDB_HOST            = "localhost"
	IMMUDB_PORT            = 3322
	IMMUDB_SSL             = false
	DB_TABLE               = "health"
	SQL_IMMUDB_DRIVER_NAME = "instrumented-sqlx-immudb"
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
		WithAddress(IMMUDB_HOST).
		WithPort(IMMUDB_PORT)

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

	db, err := sqlx.Open(
		SQL_IMMUDB_DRIVER_NAME,
		fmt.Sprintf("immudb://%s:%s@%s:%d/%s?sslmode=disable", IMMUDB_USER, IMMUDB_PWD, IMMUDB_HOST, IMMUDB_PORT, IMMUDB_DB_NAME),
	)
	handleErr(err)

	// Create our test table
	_, err = client.SQLExec(ctx, `
		CREATE TABLE IF NOT EXISTS health(
    		id INTEGER AUTO_INCREMENT,
    		name VARCHAR,
   			was_successful boolean NOT NULL,
    		PRIMARY KEY (id)
		);
	`, map[string]interface{}{})
	if err != nil {
		log.Fatal("failed to open DB", err)
	}

	// Insert Data with PGX
	sqlResult := db.MustExec(`INSERT INTO health (name, was_successful) VALUES ($1, $2)`, "Austin", true)

	// PGX gives us the the Last ID, we'll use this to verify our transaction below
	lastId, err := sqlResult.LastInsertId()
	if err != nil {
		log.Fatal(err, " - last insert ID")
	}

	// From here, we're going to construct the arguments for the client.VerifyRow method
	// We'll need
	// 1. The column names - you can get these by querying the data as shown below
	// 	queryResult, err := client.SQLQuery(ctx, `SELECT * FROM health WHERE id = @id`, map[string]interface{}{"id": lastId}, true)
	// 2. The Values we expect to be found in the column with the specified primary key, in the []*schema.SQLValue type
	// 3. The primary key stored as []*schema.SQLValue

	verifyRow := &schema.Row{
		// Here are the column names, as a reminder you can get these out of a queryResult
		Columns: []string{
			"(health.id)",
			"(health.name)",
			"(health.was_successful)",
		},
		// The values we are expecting to find
		Values: []*schema.SQLValue{
			{
				Value: &schema.SQLValue_N{
					N: lastId,
				},
			},
			{
				Value: &schema.SQLValue_S{
					S: "Austin",
				},
			},
			{
				Value: &schema.SQLValue_B{
					B: true,
				},
			},
		},
	}

	// The Primary Key
	PK := []*schema.SQLValue{
		{
			Value: &schema.SQLValue_N{
				N: lastId,
			},
		},
	}

	// Verify the row
	err = client.VerifyRow(ctx, verifyRow, DB_TABLE, PK)
	if err != nil {
		log.Fatal(err)
	}
}
