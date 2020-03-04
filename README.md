## Development

### Architecture

[clean-architecture]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

This project follows [The Clean Architecture][clean-architecture] approach. Check the following
links for **reference** only, as we don't follow them exactly:

- https://hackernoon.com/golang-clean-archithecture-efd6d7c43047
- https://medium.com/@eminetto/clean-architecture-using-golang-b63587aa5e3f
- https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1

### Folder structure

- `/internal`: Contains code that we use internally (not part of the public API).
- `/models`: Contains the base entities/structures that serve as building blocks.
- `/repositories`: Contains the interfaces for interacting with the data-store.
- `/repositories/postgres`: Contains the interfaces implementations for postgres.

### Libraries

[echo]: https://github.com/labstack/echo
[sqlx]: https://github.com/jmoiron/sqlx
[pgx]: https://github.com/JackC/pgx
[testify/assert]: https://github.com/stretchr/testify

- [echo][echo]: Web micro-framework for our REST/presentation layer.
- [sqlx][sqlx]: Database extensions for Go's standard library.
- [pgx][pgx]: Postgres database driver for Go's standard library.
- [testify/assert]: Assertions framework for our tests.

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
