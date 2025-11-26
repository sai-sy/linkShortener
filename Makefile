# -----------------------------------
# Docker Compose Helpers
# -----------------------------------

DC = docker compose

# -----------------------------------
# Database Controls
# -----------------------------------

## Start ONLY the Postgres database service
db-up:
	$(DC) up -d db

## Stop only the DB (does not delete volume)
db-down:
	$(DC) stop db

## Completely destroy the DB volume (DANGEROUS - wipes all data)
db-reset:
	$(DC) down -v
	$(DC) up -d db
	sleep 3
	$(DC) run --rm migrate

# -----------------------------------
# Migrations
# -----------------------------------

## Apply all migrations (migrate up)
migrate-up:
	$(DC) run --rm migrate

## Roll back the last migration (migrate down)
migrate-down:
	$(DC) run --rm migrate down

## Create a new migration using migrate CLI (requires migrate installed locally)
migrate-create:
	docker run -it --rm --network host --volume "$(PWD)/cmd/migrate/migrations:/migrations" migrate/migrate create -ext sql -dir /migrations "$(name)"

# -----------------------------------
# Dev Environment
# -----------------------------------

## Start full development environment (db + app)
dev-up:
	$(DC) up -d db app

## Start full dev env AND apply migrations
dev:
	$(DC) up -d db 
	$(DC) run --rm migrate
	$(DC) up app

## Stop dev environment
dev-down:
	$(DC) down

# -----------------------------------
# Utility
# -----------------------------------

## Show logs for the app
logs:
	$(DC) logs -f app

## Rebuild the dev app image
app-build:
	$(DC) build app

SQLC_IMAGE = sqlc/sqlc:1.27.0  # pick your sqlc version

sqlc:
	docker run --rm \
	  -v $(PWD):/src \
	  -w /src \
	  $(SQLC_IMAGE) generate
