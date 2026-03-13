migration:
	@migrate create -ext sql -dir database/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

run-api: 
	@go run cmd/api/main.go