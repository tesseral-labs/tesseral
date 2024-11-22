.PHONY: queries
queries:
	pg_dump --schema-only 'postgres://postgres:password@localhost?sslmode=disable' > sqlc/schema.sql
	sqlc -f ./sqlc/sqlc.yaml generate

.PHONY: proto
proto:
	npx buf generate --template buf/buf.gen.yaml
	buf generate --template buf/buf.gen-openapi-backend.yaml
	buf generate --template buf/buf.gen-openapi-frontend.yaml
	buf generate --template buf/buf.gen-openapi-intermediate.yaml
