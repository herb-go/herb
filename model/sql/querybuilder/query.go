package querybuilder

import (
	"database/sql"
)

const indexTableName = 0
const indexAlias = 1

// Query sql query interface
type Query interface {
	// QueryCommand return query command
	QueryCommand() string
	// QueryArgs return query args
	QueryArgs() []interface{}
}

// New create new plain query with given command and args
func (b *Builder) New(command string, args ...interface{}) *PlainQuery {
	return &PlainQuery{
		Builder: b,
		Command: command,
		Args:    args,
	}
}

// PlainQuery plain query struct
type PlainQuery struct {
	Builder *Builder
	// Command query command
	Command string
	// Args qyert args
	Args []interface{}
}

// QueryCommand return query command
func (q *PlainQuery) QueryCommand() string {
	if q == nil {
		return ""
	}
	return q.Command
}

// QueryArgs return query args
func (q *PlainQuery) QueryArgs() []interface{} {
	if q == nil {
		return []interface{}{}
	}
	return q.Args
}

// Exec exec query in given db
func (q *PlainQuery) Exec(db DB) (sql.Result, error) {
	return q.Builder.Exec(db, q)
}

// MustExec exec query in given db.
// Panic if any error raised.
func (q *PlainQuery) MustExec(db DB) sql.Result {
	r, err := db.Exec(q.QueryCommand(), q.QueryArgs()...)
	if err != nil {
		panic(err)
	}
	return r
}

//QueryRow query rowsfrom db with given query.
//query will be convert by builder driver.
func (q *PlainQuery) QueryRow(db DB) *sql.Row {
	return q.Builder.QueryRow(db, q)

}

//QueryRows query rows from db with given query.
//query will be convert by builder driver.
func (q *PlainQuery) QueryRows(db DB) (*sql.Rows, error) {
	return q.Builder.QueryRows(db, q)
}

// And concat query with given query list by AND operation
func (q *PlainQuery) And(qs ...Query) *PlainQuery {
	if q != nil && q.Command != "" {
		qslice := make([]Query, len(qs)+1)
		qslice[0] = q
		copy(qslice[1:], qs)
		*q = *(q.Builder.And(qslice...))
	} else {
		*q = *(q.Builder.And(qs...))
	}
	return q
}

// Or concat query with given query list by OR operation
func (q *PlainQuery) Or(qs ...Query) *PlainQuery {
	if q != nil && q.Command != "" {
		qslice := make([]Query, len(qs)+1)
		qslice[0] = q
		copy(qslice[1:], qs)
		*q = *(q.Builder.Or(qslice...))
	} else {
		*q = *(q.Builder.Or(qs...))
	}
	return q
}

// QueryData query data struct
type QueryData struct {
	// Field data field
	Field string
	// Data data value
	Data []interface{}
	// Raw data raw value.
	// if raw is setted,Data field will be ignored.
	Raw string
}
