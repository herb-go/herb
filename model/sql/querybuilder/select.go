package querybuilder

import (
	"database/sql"
)

// NewSelectClause create new select clause
func (b *Builder) NewSelectClause() *SelectClause {
	return &SelectClause{
		Builder:   b,
		Prefix:    b.New(""),
		Fields:    []string{},
		fieldargs: []interface{}{},
	}
}

// SelectClause select clause struct
type SelectClause struct {
	Builder   *Builder
	Prefix    *PlainQuery
	Fields    []string
	fieldargs []interface{}
}

// AddFields add fields to select clause
func (q *SelectClause) AddFields(m *Fields) *SelectClause {
	var fields = make([]string, len(*m))
	var i = 0
	for k := range *m {
		fields[i] = (*m)[k].Field
		i++
	}
	return q.Add(fields...)
}

// Add add field to select clause
func (q *SelectClause) Add(fields ...string) *SelectClause {
	q.Fields = append(q.Fields, fields...)
	return q
}

// AddRaw add raw fields to select clause
func (q *SelectClause) AddRaw(fields ...interface{}) *SelectClause {
	for k := range fields {
		q.Fields = append(q.Fields, "?")
		q.fieldargs = append(q.fieldargs, fields[k])
	}
	return q
}

// AddSelect add select subquery to select clause
func (q *SelectClause) AddSelect(Select *SelectQuery) *SelectClause {
	query := *Select.Query()
	q.Fields = append(q.Fields, "("+query.QueryCommand()+")")
	q.fieldargs = append(q.fieldargs, query.QueryArgs()...)

	return q
}

// QueryCommand return query command
func (q *SelectClause) QueryCommand() string {
	var command = "SELECT"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	var columns = " "
	for k := range q.Fields {
		columns += q.Fields[k] + " , "
	}
	if len(q.Fields) > 0 {
		columns = columns[:len(columns)-3]
	}
	command += columns
	return command
}

// QueryArgs return query args
func (q *SelectClause) QueryArgs() []interface{} {
	args := []interface{}{}
	args = append(args, q.Prefix.QueryArgs()...)
	args = append(args, q.fieldargs...)
	return args
}

// Result return select result with select clause
func (q *SelectClause) Result() *SelectResult {
	return NewSelectResult(q.Fields)
}

// NewSelectResult create select result with given fields
func NewSelectResult(fields []string) *SelectResult {
	return &SelectResult{
		Fields: fields,
		args:   make([]interface{}, len(fields)),
	}

}

// ResultScanner select result scanner interface
type ResultScanner interface {
	Scan(dest ...interface{}) error
}

// SelectResult select result struct
type SelectResult struct {
	Fields []string
	args   []interface{}
}

// Bind bind field and value pointer to select result.
func (r *SelectResult) Bind(field string, pointer interface{}) *SelectResult {
	for k := range r.Fields {
		if r.Fields[k] == field {
			r.args[k] = pointer
			return r
		}
	}
	return r
}

// BindFields bind fields to select result
func (r *SelectResult) BindFields(m *Fields) *SelectResult {
	for _, v := range *m {
		r.Bind(v.Field, v.Data)
	}
	return r
}

// Pointers return field pointers
func (r *SelectResult) Pointers() []interface{} {
	return r.args
}

//ScanFrom scan data with result scanner
func (r *SelectResult) ScanFrom(s ResultScanner) error {
	return s.Scan(r.Pointers()...)
}

// NewSelectQuery create new select
func (b *Builder) NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Builder: b,
		Select:  b.NewSelectClause(),
		From:    b.NewFromClause(),
		Join:    b.NewJoinClause(),
		Where:   b.NewWhereClause(),
		Having:  b.NewHavingClause(),
		GroupBy: b.NewGroupByClause(),
		OrderBy: b.NewOrderByClause(),
		Limit:   b.NewLimitClause(),
		Other:   b.New(""),
	}
}

// SelectQuery select query struct.
type SelectQuery struct {
	Builder *Builder
	Select  *SelectClause
	From    *FromClause
	Join    *JoinClause
	Where   *WhereClause
	GroupBy *GroupByClause
	Having  *HavingClause
	OrderBy *OrderByClause
	Limit   *LimitClause
	Other   *PlainQuery
}

// Result return select result
func (s *SelectQuery) Result() *SelectResult {
	return s.Select.Result()
}

// Query convert select query to plain query.
func (s *SelectQuery) Query() *PlainQuery {
	return s.Builder.Lines(s.Select, s.From, s.Join, s.Where, s.GroupBy, s.Having, s.OrderBy, s.Limit, s.Other)
}

// QueryCommand return query command
func (s *SelectQuery) QueryCommand() string {
	return s.Query().Command
}

// QueryArgs return query args
func (s *SelectQuery) QueryArgs() []interface{} {
	return s.Query().Args
}

//QueryRow query rowsfrom db with given query.
//query will be convert by builder driver.
func (s *SelectQuery) QueryRow(db DB) *sql.Row {
	return s.Builder.QueryRow(db, s)

}

//QueryRows query rows from db with given query.
//query will be convert by builder driver.
func (s *SelectQuery) QueryRows(db DB) (*sql.Rows, error) {
	return s.Builder.QueryRows(db, s)
}
