package querybuilder

import (
	"database/sql"
)

const indexTableName = 0
const indexAlias = 1

type Query interface {
	QueryCommand() string
	QueryArgs() []interface{}
}

func (b *Builder) New(command string, args ...interface{}) *PlainQuery {
	return &PlainQuery{
		Builder: b,
		Command: command,
		Args:    args,
	}
}

type PlainQuery struct {
	Builder *Builder
	Command string
	Args    []interface{}
}

func (q *PlainQuery) QueryCommand() string {
	if q == nil {
		return ""
	}
	return q.Command
}

func (q *PlainQuery) QueryArgs() []interface{} {
	if q == nil {
		return []interface{}{}
	}
	return q.Args
}
func (q *PlainQuery) Exec(db DB) (sql.Result, error) {
	return q.Builder.Exec(db, q)
}
func (q *PlainQuery) MustExec(db DB) sql.Result {
	r, err := db.Exec(q.QueryCommand(), q.QueryArgs()...)
	if err != nil {
		panic(err)
	}
	return r
}

func (q *PlainQuery) QueryRow(db DB) *sql.Row {
	return q.Builder.QueryRow(db, q)

}
func (q *PlainQuery) QueryRows(db DB) (*sql.Rows, error) {
	return q.Builder.QueryRows(db, q)
}
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

type QueryData struct {
	Field string
	Data  []interface{}
	Raw   string
}
