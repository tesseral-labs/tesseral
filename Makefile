.PHONY: dev
dev:
	docker-compose up --build --watch

.PHONY: queries
queries:
	pg_dump --schema-only 'postgres://postgres:password@localhost?sslmode=disable' > sqlc/schema.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	sqlc -f ./sqlc/sqlc.yaml generate

.PHONY: proto
proto:
	npx buf generate --template buf/buf.gen.yaml
	buf generate --template buf/buf.gen-openapi-backend.yaml
	buf generate --template buf/buf.gen-openapi-frontend.yaml
	buf generate --template buf/buf.gen-openapi-intermediate.yaml
