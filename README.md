# immutable-databases

Immutable databases performance and features analysis

## How to configure environment and run

1. Run `docker compose up -d` to start the MinIO and immudb container
1. Open http://localhost:9001/login to login into the MinIO console (user: `minio-user`, password: `minio-password`)
1. Open http://localhost:8080 to access the immudb web console (user: `immudb`, password: `immudb`)
1. Run `go mod tidy` to download dependencies
1. Run one of the following commands to run the application

```bash
# Simple immudb client
go run cmd/immudb/main.go
```

or

```bash
# SQLX client with logging
go run cmd/immudb-sqlx-instrumented/main.go
```

or

```bash
# SQL standard client with logging
go run cmd/immudb-stdlib-instrumented/main.go
```

or

```bash
# SQL standard client with transactions
go run cmd/immudb-stdlib-tx/main.go
```

## Resources

- [immudb](https://github.com/codenotary/immudb)
- [immudb SQL transactions](https://docs.immudb.io/master/develop/sql/transactions.html)
- [immudb S3 storage](https://docs.immudb.io/master/production/s3-storage.html)
