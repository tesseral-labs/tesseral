.PHONY: bootstrap
bootstrap:
	@# Install the required ssl certs
	./bin/create-localhost-certs
	@# Start the docker containers
	docker compose up -d --wait
	@# Wait for the database to be ready
	@until PGPASSWORD=password psql "postgres://postgres:password@localhost:5432?sslmode=disable" -c "SELECT 1" >/dev/null 2>&1; do \
		echo "PostgreSQL is unavailable - retrying..."; \
		sleep 2; \
	done
	@# Run database migrations
	make migrate up
	@# Seed the database
	psql "postgres://postgres:password@localhost:5432?sslmode=disable" -f .local/db/seed.sql
	@# Stop the docker containers
	docker compose stop

.PHONY: dev
dev:
	docker compose up --build --watch

.PHONY: migrate
ARGS = $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
migrate:
	migrate -path cmd/openauthctl/migrations -database "postgres://postgres:password@localhost:5432?sslmode=disable" $(ARGS)
%:
	@:

.PHONY: proto
proto:
	rm -rf internal/backend/gen internal/frontend/gen internal/intermediate/gen internal/oauth/gen internal/common/gen console/src/gen ui/src/gen
	npx buf generate --template buf/buf.gen-backend.yaml
	npx buf generate --template buf/buf.gen-frontend.yaml
	npx buf generate --template buf/buf.gen-intermediate.yaml
	npx buf generate --template buf/buf.gen-oauth.yaml
	npx buf generate --template buf/buf.gen-common.yaml

.PHONY: queries
queries:
	rm -rf internal/store/queries internal/backend/store/queries internal/frontend/store/queries internal/intermediate/store/queries internal/oauth/store/queries internal/saml/store/queries internal/scim/store/queries internal/common/store/queries internal/wellknown/store/queries
	docker run --rm --volume "$$(pwd)/sqlc/queries.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-backend.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-frontend.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-intermediate.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-oauth.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-saml.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-scim.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-common.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-wellknown.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	sqlc -f ./sqlc/sqlc.yaml generate
