#!/bin/bash

/go/bin/migrate -path database/migrations/ -database "$POSTGRES_DB_URL" -verbose up