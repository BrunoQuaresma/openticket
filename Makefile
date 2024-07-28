all: ./api/database/sqlc.yaml ./web/src/sdk/types.gen.ts ./web/pnpm-lock.yaml

./api/database/sqlc.yaml: ./api/database/queries/*.sql
	sqlc generate --file ./api/database/sqlc.yaml

./web/src/sdk/types.gen.ts: ./api/*.go ./tygo.yaml
	tygo generate

./web/pnpm-lock.yaml: ./web/package.json
	cd ./web && pnpm install