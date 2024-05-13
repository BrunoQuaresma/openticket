gen: ./api/database/queries/*.sql
	sqlc generate --file ./api/database/sqlc.yaml