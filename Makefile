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
	@# Create the tmp directory for ecdsa keys
	mkdir -p .local/tmp
	@# Generate the session signing key
	openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out .local/tmp/session-signing-key.pem
	sed -e '1d' -e '$$d' .local/tmp/session-signing-key.pem > .local/tmp/trimmed-session-signing-key.pem
	openssl ec -in .local/tmp/session-signing-key.pem -pubout -out .local/tmp/session-signing-public-key.pem
	@# Encrypt the session signing key with KMS
	AWS_DEFAULT_REGION=us-west-1 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test aws kms encrypt --encryption-algorithm "RSAES_OAEP_SHA_256" --endpoint-url "http://localhost:4566" --key-id "bc436485-5092-42b8-92a3-0aa8b93536dc" --output text --plaintext fileb://.local/tmp/trimmed-session-signing-key.pem --query CiphertextBlob | base64 -d > .local/tmp/session-signing-key.encrypted
	@# Seed the database
	psql "postgres://postgres:password@localhost:5432?sslmode=disable" -f .local/db/seed.sql
	@# Stop the docker containers
	docker compose stop
	@# Remove the tmp directory
	rm -rf .local/tmp

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
	rm -rf internal/backend/gen internal/frontend/gen internal/intermediate/gen internal/oauth/gen internal/common/gen app/src/gen ui/src/gen
	npx buf generate --template buf/buf.gen-backend.yaml
	npx buf generate --template buf/buf.gen-frontend.yaml
	npx buf generate --template buf/buf.gen-intermediate.yaml
	npx buf generate --template buf/buf.gen-oauth.yaml
	npx buf generate --template buf/buf.gen-common.yaml

.PHONY: queries
queries:
	rm -rf internal/store/queries internal/backend/store/queries internal/frontend/store/queries internal/intermediate/store/queries internal/oauth/store/queries internal/saml/store/queries internal/scim/store/queries internal/common/store/queries
	docker run --rm --volume "$$(pwd)/sqlc/queries.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-backend.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-frontend.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-intermediate.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-oauth.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-saml.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-scim.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-common.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	sqlc -f ./sqlc/sqlc.yaml generate
