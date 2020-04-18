## Note to professor and TA(s)

Here's a list of folders in which you will find what you need:

1. **Schema**: Can be found at `/ddl/schema.sql`.

2. **Triggers**: Can all be found at `/ddl/triggers.sql`. Included triggers take
care of listening and emitting events for maintenance status changes, and bad
reviews.

3. **Queries for CRUD operations**: Check the non-test files at `/repositories/postgres/*.go`.
For example, `/repositories/postgres/ride.go` contains all the CRUD operations for rides.

4. **Queries for reports**: Check the files at `/reports/*.read.sql`.

----

## Development

### Architecture

[clean-architecture]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

This project follows [The Clean Architecture][clean-architecture] approach. Check the following
links for **reference** only, as we don't follow them exactly:

- https://hackernoon.com/golang-clean-archithecture-efd6d7c43047
- https://medium.com/@eminetto/clean-architecture-using-golang-b63587aa5e3f
- https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1

### Folder structure

**Note**: Folders that are not too important were skipped.

- `/generator`: Containst the script and functions for generator mock data.

- `/handlers`: Contains the HTTP handlers for echo.

- `/internal`: Contains code that we use internally (not part of the public API).

- `/models`: Contains the base entities/structures that serve as the building
blocks of our application.

- `/reports`: Contains the SQL files that generate our reports.

- `/repositories`: Contains the interfaces for interacting with the data-store.
Note that these interfaces define how to *Create*, *Read*, *Update*, and
*Delete* (CRUD) **models**; they **do not** implement business-logic or data
validation, as that's the job of the **usecases**.

- `/repositories/postgres`: Contains the interfaces implementations for postgres.

- `/usecases`: Contains the interfaces and implementation of our business-logic.
**usecases** are the heart of our application, as they take care of taking the
input, and using the appropriate models and repositories to execute a business
rule. For example:

  - Creating a new user should validate that their username or email doesn't
  exists already. We do this by getting the parameters of the use case.

  - Changing the email of the user validates that the requested email is a
  valid email address.

  - Changing the gender of the user should validate that the requested gender
  is actually valid.

  - etc.

### Libraries

[echo]: https://github.com/labstack/echo
[sqlx]: https://github.com/jmoiron/sqlx
[pgx]: https://github.com/JackC/pgx
[squirrel]: https://github.com/Masterminds/squirrel
[testify/assert]: https://github.com/stretchr/testify
[prest]: https://github.com/prest/prest

- [echo][echo]: Web micro-framework for our REST/presentation layer.
- [sqlx][sqlx]: Database extensions for Go's standard library.
- [pgx][pgx]: Postgres database driver for Go's standard library.
- [squirrel][squirrel]: Re-usable SQL builder.
- [testify/assert]: Assertions framework for our tests.
- [prest]: App that takes care of serving our reports.

### Testing

If you want to run the backend locally you will need to go though the database
setup, after which you need to start the HTTP server.

#### Database teardown

Before you setup the database you might want to destroy a previous instance;
to do so just run:

```sh
docker-compose down
```

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

**Third**, populate the database with mock that:

```sh
go run main.go generate
```

**Fourth**, create triggers:

```sh
docker-compose exec postgres psql -U postgres -d testdb -f /ddl/triggers.sql
```

#### Running tests

On the root folder of this project:

```sh
go test ./...
```

#### Populating database

```sh
go run main.go generate
```

#### Running HTTP server

In dev environment:

```sh
go run main.go -t server
```

In production environment:

```sh
go run main.go server
```
