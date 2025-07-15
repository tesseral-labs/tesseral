.PHONY: bootstrap
bootstrap:
	@# Shut down and clear out any local postgres state
	docker compose stop postgres
	rm -rf .local/postgres
	@# Install the required ssl certs
	./bin/create-localhost-certs
	@# Start the database docker container
	docker compose up -d --wait postgres
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
	docker compose stop postgres

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
	rm -rf internal/auditlog/gen internal/backend/gen internal/frontend/gen internal/intermediate/gen internal/common/gen console/src/gen vault-ui/src/gen
	buf format internal/auditlog/proto -w
	buf format internal/backend/proto -w
	buf format internal/frontend/proto -w
	buf format internal/intermediate/proto -w
	buf format internal/common/proto -w
	npx buf generate --template buf/buf.gen-auditlog.yaml
	npx buf generate --template buf/buf.gen-backend.yaml
	npx buf generate --template buf/buf.gen-frontend.yaml
	npx buf generate --template buf/buf.gen-intermediate.yaml
	npx buf generate --template buf/buf.gen-common.yaml

.PHONY: queries
queries:
	rm -rf \
		internal/store/queries \
		internal/auditlog/store/queries \
		internal/backend/store/queries \
		internal/frontend/store/queries \
		internal/intermediate/store/queries \
		internal/saml/store/queries \
		internal/oidc/store/queries \
		internal/scim/store/queries \
		internal/common/store/queries \
		internal/configapi/store/queries \
		internal/defaultoauth/store/queries
	./bin/pg_format/pg_format -i sqlc/queries.sql
	./bin/pg_format/pg_format -i sqlc/queries-auditlog.sql
	./bin/pg_format/pg_format -i sqlc/queries-backend.sql
	./bin/pg_format/pg_format -i sqlc/queries-frontend.sql
	./bin/pg_format/pg_format -i sqlc/queries-intermediate.sql
	./bin/pg_format/pg_format -i sqlc/queries-saml.sql
	./bin/pg_format/pg_format -i sqlc/queries-oidc.sql
	./bin/pg_format/pg_format -i sqlc/queries-scim.sql
	./bin/pg_format/pg_format -i sqlc/queries-common.sql
	./bin/pg_format/pg_format -i sqlc/queries-configapi.sql
	./bin/pg_format/pg_format -i sqlc/queries-defaultoauth.sql
	sqlc -f ./sqlc/sqlc.yaml generate
