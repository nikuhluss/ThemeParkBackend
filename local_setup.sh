#!/usr/bin/env bash

docker-compose exec postgres psql -U postgres -c'DROP DATABASE IF EXISTS testdb'
docker-compose exec postgres psql -U postgres -c'CREATE DATABASE testdb'
docker-compose exec postgres psql -U postgres -d testdb -c'CREATE SCHEMA theme_park'
docker-compose exec postgres psql -U postgres -d testdb -f /ddl/schema.sql

go run main.go generate

docker-compose exec postgres psql -U postgres -d testdb -f /ddl/triggers.sql
