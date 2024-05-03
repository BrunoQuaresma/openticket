#!/bin/bash

migrate -path database/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up