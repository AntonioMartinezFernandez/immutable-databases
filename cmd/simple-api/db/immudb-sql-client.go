package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/config"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
)

func NewImmuDbSqlClient(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	// Stablish database connection configuration
	opts := immudb.DefaultOptions().
		WithAddress(cfg.ImmudbHost).
		WithPort(cfg.ImmudbPort).
		WithUsername(cfg.ImmudbUser).
		WithPassword(cfg.ImmudbPwd).
		WithDatabase(cfg.ImmudbDatabase)

	// Create a new standard database/sql client and stablish pool connection
	dbClient := stdlib.OpenDB(opts)
	dbClient.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeInMinutes) * time.Minute)
	dbClient.SetMaxOpenConns(cfg.MaxConnections)
	dbClient.SetMaxIdleConns(cfg.ConnIdle)

	pingErr := dbClient.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	// Create table if not exists
	dbCreateSqlQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR[128],
			streamId VARCHAR[128],
			content VARCHAR[65535],
			PRIMARY KEY id
			);`,
		cfg.DbTable,
	)

	_, ctErr := dbClient.ExecContext(ctx, dbCreateSqlQuery)
	if ctErr != nil {
		return nil, ctErr
	}

	return dbClient, nil
}
