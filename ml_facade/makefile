.PHONY: app/run
app/run:
	@go run ./cmd

.PHONY: app/run-dev
app/run-dev:
	@go run ./cmd -db-host localhost -rdb-host localhost -ml-service-host localhost

.PHONY: app/dev-containers
app/run-dev-containers:
	@docker compose up -d

.PHONY: app/test
app/test:
	@go test ./...

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${MONITORING_DB_DSN} up

.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
