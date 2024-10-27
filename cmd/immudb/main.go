package main

import (
	"context"
	"fmt"
	"log"
	"time"

	immudb "github.com/codenotary/immudb/pkg/client"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	opts := immudb.DefaultOptions().
		WithAddress("localhost").
		WithPort(3322)

	client := immudb.NewClient().WithOptions(opts)

	err := client.OpenSession(
		ctx,
		[]byte(`immudb`),
		[]byte(`immudb`),
		"defaultdb",
	)
	handleErr(err)

	defer client.CloseSession(ctx)

	tx, err := client.NewTx(ctx)
	handleErr(err)

	err = tx.SQLExec(
		ctx,
		`CREATE TABLE IF NOT EXISTS posts(id INTEGER AUTO_INCREMENT, personId VARCHAR[128], text VARCHAR[4096], active BOOLEAN, PRIMARY KEY id);`,
		nil,
	)
	handleErr(err)

	txh, err := tx.Commit(ctx)
	handleErr(err)

	fmt.Printf("Successfully committed rows %d creating %s table\n", txh.UpdatedRows, "posts")

	start := time.Now()

	nRows := 1000
	for i := 0; i < nRows; i++ {
		txRows, err := client.NewTx(ctx)
		handleErr(err)
		sqlExecErr := txRows.SQLExec(
			ctx,
			"INSERT INTO posts(personId, text, active) VALUES (@personId,@text, @active)",
			map[string]interface{}{
				"personId": uuid.New().String(),
				"text":     faker.Paragraph(),
				"active":   i%2 == 0,
			},
		)
		handleErr(sqlExecErr)
		txhRows, err := txRows.Commit(ctx)
		handleErr(err)
		fmt.Printf("Successfully committed rows %d\n", txhRows.UpdatedRows)
	}

	reader, err := client.SQLQueryReader(ctx, "SELECT * FROM posts", nil)
	handleErr(err)

	for reader.Next() {
		row, err := reader.Read()
		handleErr(err)

		fmt.Println(row[0], row[1])
	}

	fmt.Println("Time taken to insert 1000 rows one by one: ", time.Since(start))
}
