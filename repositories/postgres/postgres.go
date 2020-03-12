package postgres

import (
	sq "github.com/Masterminds/squirrel"
)

// psql is a statement builder that uses Dollar format (Postgres).
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
