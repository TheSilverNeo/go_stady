ARGS=$(filter-out $@, $(MAKECMDGOALS))

DB_NAME=$$APP_DB_NAME
POSTGRESQL_URL=postgresql://$$APP_DB_USER:$$APP_DB_PASSWORD@$$APP_DB_HOST:$$APP_DB_PORT/$(DB_NAME)?sslmode=$$APP_DB_SSL_MODE
MGR_COMMAND=-path=$$APP_DB_MIGRATIONS_DIR -database $(POSTGRESQL_URL) -verbose

init: create-env \
 	  migrate-up

create-env:
	cp --update=none .env.sample .env

migrate-create:
	export `grep -v "^#" .env | xargs` && \
	docker compose run --rm migrate $(MGR_COMMAND) create -ext sql -dir $$APP_DB_MIGRATIONS_DIR -seq $(ARGS)
	sudo chown -R $$USER. ./migrations

migrate-up:
	export `grep -v "^#" .env | xargs` && \
	docker compose run --rm migrate $(MGR_COMMAND) up

migrate-down:
	export `grep -v "^#" .env | xargs` && \
	env docker-compose run --rm migrate $(MGR_COMMAND) down $(ARGS)