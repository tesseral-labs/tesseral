.PHONY: dev
dev:
	docker compose up --build --watch

.PHONY: queries
queries:
	rm -r internal/store/queries internal/backend/store/queries internal/frontend/store/queries internal/intermediate/store/queries internal/oauth/store/queries internal/saml/store/queries internal/scim/store/queries internal/common/store/queries
	docker run --rm --volume "$$(pwd)/sqlc/queries.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-backend.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-frontend.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-intermediate.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-oauth.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-saml.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-scim.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	docker run --rm --volume "$$(pwd)/sqlc/queries-shared.sql:/work/queries.sql" backplane/pgformatter -i queries.sql
	sqlc -f ./sqlc/sqlc.yaml generate

.PHONY: proto
proto:
	rm -r internal/backend/gen internal/frontend/gen internal/intermediate/gen internal/oauth/gen internal/common/gen app/src/gen ui/src/gen
	npx buf generate --template buf/buf.gen-backend.yaml
	npx buf generate --template buf/buf.gen-frontend.yaml
	npx buf generate --template buf/buf.gen-intermediate.yaml
	npx buf generate --template buf/buf.gen-oauth.yaml
	npx buf generate --template buf/buf.gen-common.yaml
