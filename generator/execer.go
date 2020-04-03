package generator

import "database/sql"

// Execer interface represents an object with the Exec function.
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// MustExecer interface represents an object with the MustExec function.
type MustExecer interface {
	MustExec(query string, args ...interface{}) sql.Result
}

// ExecerMustExecer combines both Execer and MustExecer.
type ExecerMustExecer interface {
	Execer
	MustExecer
}

// AsExecer takes the given MustExecer act as an Execer (error always returns nil).
type AsExecer struct {
	mustExecer MustExecer
}

// Exec executes the given query string. Error always returns nil, and panics if
// the internal MustExecer fails.
func (ae *AsExecer) Exec(query string, args ...interface{}) (sql.Result, error) {
	res := ae.mustExecer.MustExec(query, args...)
	return res, nil
}

// MustInsert is a function used to wrap Insert* calls that return an ID and
// an error. E.g.
// ```
// MustInsert(InsertSomething())
// ```
func MustInsert(returnedID string, returnedError error) string {
	if returnedError != nil {
		panic(returnedError)
	}
	return returnedID
}
