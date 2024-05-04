MIGRATIONS_DIR=./internal/remote/db/migrations

new-migration:
	@read -p "New migration filename: " filename;\
		goose -dir $(MIGRATIONS_DIR) create $$filename go

start-dev-psql:
	@docker run -d --rm --name dev-psql -e POSTGRES_PASSWORD=password -e POSTGRES_DB=postgres -p 5432:5432 postgres

stop-dev-psql:
	@docker stop dev-psql
