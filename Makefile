all: ./api/database/sqlc.yaml ./web/src/sdk/types.gen.ts

./api/database/sqlc.yaml: ./api/database/queries/*.sql
	sqlc generate --file ./api/database/sqlc.yaml

./web/src/sdk/types.gen.ts: ./api/*.go ./tygo.yaml
	tygo generate