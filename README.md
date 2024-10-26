# immutable-databases

Immutable databases performance and features analysis

## How to configure environment and run

```bash
docker compose up -d
```

1. Open http://localhost:9001/login to login into the MinIO console (user: `minio-user`, password: `minio-password`)
1. Create a bucket called `immudb-bucket`
1. Open http://localhost:8080 to access the immudb web console (user: `immudb`, password: `immudb`)

```bash
go run cmd/main.go
```

## Resources

- [immudb](https://github.com/codenotary/immudb)
- [immudb SQL transactions](https://docs.immudb.io/master/develop/sql/transactions.html)
- [immudb S3 storage](https://docs.immudb.io/master/production/s3-storage.html)
