# Run Load tests with k6
test:
  @k6 run tests/k6-start-event.js

# Run Load tests with vegeta and export report
vegeta:
  @cd tests/vegeta && go run ./main.go
  @echo "Plot it with 'https://hdrhistogram.github.io/HdrHistogram/plotFiles.html'"

# Start necessary services to run the integration tests or run the application in development mode
infra:
  @docker compose up -d

# Run app locally as independent service
run:
  @go mod tidy
  @go run cmd/simple-api/main.go

