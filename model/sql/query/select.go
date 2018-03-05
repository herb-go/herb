package query

import (
	"database/sql"
	"time"
)

func (b *Builder) NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Builder: b,
		Prefix:  b.New(""),
		Fields:  []string{},
	}
}

type SelectQuery struct {
	Builder *Builder
	Prefix  *PlainQuery
	Fields  []string
}

func (q *SelectQuery) AddFields(m Fields) *SelectQuery {
	var fields = make([]string, len(m))
	var i = 0
	for k := range m {
		fields[i] = k
		i++
	}
	return q.Add(fields...)
}
func (q *SelectQuery) Add(fields ...string) *SelectQuery {
	q.Fields = append(q.Fields, fields...)
	return q
}

func (q *SelectQuery) QueryCommand() string {
	var command = "SELECT"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	var colunms = " "
	for k := range q.Fields {
		colunms += q.Fields[k] + " , "
	}
	command += colunms[:len(colunms)-3]
	return command
}
func (q *SelectQuery) QueryArgs() []interface{} {
	return q.Prefix.QueryArgs()
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

func (r *SelectResult) BindFields(m Fields) *SelectResult {
	for k, v := range m {
		r.Bind(k, v)
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
		Other:   b.New(""),
	}
}

type Select struct {
	Builder *Builder
	Select  *SelectQuery
	From    *FromQuery
	Join    *JoinQuery
	Where   *WhereQurey
	Other   *PlainQuery
}

func (s *Select) Result() *SelectResult {
	return s.Select.Result()
}

func (s *Select) Query() *PlainQuery {
	return s.Builder.Lines(s.Select, s.From, s.Join, s.Where, s.Other)
}
func (s *Select) QueryRow(db DB) *sql.Row {
	q := s.Query()
	cmd := q.QueryCommand()
	args := q.QueryArgs()
	var timestamp int64
	if Debug {
		timestamp = time.Now().UnixNano()
	}
	err := db.QueryRow(cmd, args...)
	if Debug {
		Logger(timestamp, cmd, args)
	}
	return err
}
func (s *Select) QueryRows(db DB) (*sql.Rows, error) {
	q := s.Query()
	cmd := q.QueryCommand()
	args := q.QueryArgs()
	var timestamp int64
	if Debug {
		timestamp = time.Now().UnixNano()
	}
	rows, err := db.Query(cmd, args...)
	if Debug {
		Logger(timestamp, cmd, args)
	}
	return rows, err
}
