## Development

### Testing

#### Database setup

**First**, you need to have docker installed. Once that's done, you need to run:

```sh
docker-compose up --build
```

**Second**, after the postgres database is running, execute:

```sh
docker-compose exec postgres psql -U postgres -c'CREATE DATABASE testdb'
docker-compose exec postgres psql -U postgres -d testdb -c'CREATE SCHEMA theme_park'
docker-compose exec postgres psql -U postgres -d testdb -f /ddl/schema.sql
```

#### Running tests

On the root folder of this project:

```sh
go test ./...
```
