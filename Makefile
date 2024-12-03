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
	npx buf generate --template buf/buf.gen-backend.yaml
	npx buf generate --template buf/buf.gen-frontend.yaml
	npx buf generate --template buf/buf.gen-intermediate.yaml
