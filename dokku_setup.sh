#!/usr/bin/env bash

if [ -z "$DATABASE_URL" ]; then
    echo "DATABASE_URL environment variable is required"
    exit 1
fi

psql="psql $DATABASE_URL"

# see: https://gist.github.com/TheMengzor/968e5ea87e99d9c41782
CURRDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
DDLDIR="$CURRDIR/ddl"

docker run -it --rm postgres $psql -c'DROP SCHEMA IF EXISTS theme_park'
docker run -it --rm postgres $psql -c'CREATE SCHEMA theme_park'
docker run -it --rm -v "$DDLDIR:/ddl" postgres $psql -f /ddl/schema.sql

go run main.go --dokku generate

docker run -it --rm -v "$DDLDIR:/ddl" postgres $psql -f /ddl/triggers.sql
