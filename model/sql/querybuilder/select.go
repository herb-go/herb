package querybuilder

import (
	"database/sql"
)

func (b *Builder) NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Builder:   b,
		Prefix:    b.New(""),
		Fields:    []string{},
		fieldargs: []interface{}{},
	}
}

type SelectQuery struct {
	Builder   *Builder
	Prefix    *PlainQuery
	Fields    []string
	fieldargs []interface{}
}

func (q *SelectQuery) AddFields(m *Fields) *SelectQuery {
	var fields = make([]string, len(*m))
	var i = 0
	for k := range *m {
		fields[i] = (*m)[k].Field
		i++
	}
	return q.Add(fields...)
}
func (q *SelectQuery) Add(fields ...string) *SelectQuery {
	q.Fields = append(q.Fields, fields...)
	return q
}

func (q *SelectQuery) AddRaw(fields ...interface{}) *SelectQuery {
	for k := range fields {
		q.Fields = append(q.Fields, "?")
		q.fieldargs = append(q.fieldargs, fields[k])
	}
	return q
}
func (q *SelectQuery) AddSelect(Select *Select) *SelectQuery {
	query := *Select.Query()
	q.Fields = append(q.Fields, "("+query.QueryCommand()+")")
	q.fieldargs = append(q.fieldargs, query.QueryArgs()...)

	return q
}

func (q *SelectQuery) QueryCommand() string {
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
func (q *SelectQuery) QueryArgs() []interface{} {
	args := []interface{}{}
	args = append(args, q.Prefix.QueryArgs()...)
	args = append(args, q.fieldargs...)
	return args
}
func (q *SelectQuery) Result() *SelectResult {
	return NewSelectResult(q.Fields)
}

func NewSelectResult(fields []string) *SelectResult {
	return &SelectResult{
		Fields: fields,
		args:   make([]interface{}, len(fields)),
	}

}

type ResultScanner interface {
	Scan(dest ...interface{}) error
}
type SelectResult struct {
	Fields []string
	args   []interface{}
}

func (r *SelectResult) Bind(field string, arg interface{}) *SelectResult {
	for k := range r.Fields {
		if r.Fields[k] == field {
			r.args[k] = arg
			return r
		}
	}
	return r
}

func (r *SelectResult) BindFields(m *Fields) *SelectResult {
	for _, v := range *m {
		r.Bind(v.Field, v.Data)
	}
	return r
}

func (r *SelectResult) Args() []interface{} {
	return r.args
}

func (r *SelectResult) ScanFrom(s ResultScanner) error {
	return s.Scan(r.Args()...)
}

func (b *Builder) NewSelect() *Select {
	return &Select{
		Builder: b,
		Select:  b.NewSelectQuery(),
		From:    b.NewFromQuery(),
		Join:    b.NewJoinQuery(),
		Where:   b.NewWhereQuery(),
		OrderBy: b.NewOrderByQuery(),
		Limit:   b.NewLimitQuery(),
		Other:   b.New(""),
	}
}

type Select struct {
	Builder *Builder
	Select  *SelectQuery
	From    *FromQuery
	Join    *JoinQuery
	Where   *WhereQuery
	OrderBy *OrderByQuery
	Limit   *LimitQuery
	Other   *PlainQuery
}

func (s *Select) Result() *SelectResult {
	return s.Select.Result()
}

func (s *Select) Query() *PlainQuery {
	return s.Builder.Lines(s.Select, s.From, s.Join, s.Where, s.OrderBy, s.Limit, s.Other)
}

func (s *Select) QueryCommand() string {
	return s.Query().Command
}
func (s *Select) QueryArgs() []interface{} {
	return s.Query().Args
}

func (s *Select) QueryRow(db DB) *sql.Row {
	return s.Builder.QueryRow(db, s)

}
func (s *Select) QueryRows(db DB) (*sql.Rows, error) {
	return s.Builder.QueryRows(db, s)
}
