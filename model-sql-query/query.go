package query

import (
	"database/sql"
	"strings"
	"time"
)

const indexTableName = 0
const indexAlias = 1

type Query interface {
	QueryCommand() string
	QueryArgs() []interface{}
}

func New(command string, args ...interface{}) *PlainQuery {
	return &PlainQuery{
		Command: command,
		Args:    args,
	}
}

type PlainQuery struct {
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
	cmd := q.QueryCommand()
	args := q.QueryArgs()
	var timestamp int64
	if Debug {
		timestamp = time.Now().UnixNano()
	}
	r, err := db.Exec(cmd, args...)
	if Debug {
		Logger(timestamp, cmd, args)
	}
	return r, err
}
func (q *PlainQuery) MustExec(db DB) sql.Result {
	r, err := db.Exec(q.QueryCommand(), q.QueryArgs()...)
	if err != nil {
		panic(err)
	}
	return r
}
func (q *PlainQuery) And(qs ...Query) *PlainQuery {
	if q != nil && q.Command != "" {
		qslice := make([]Query, len(qs)+1)
		qslice[0] = q
		copy(qslice[1:], qs)
		*q = *(And(qslice...))
	} else {
		*q = *(And(qs...))
	}
	return q
}

func (q *PlainQuery) Or(qs ...Query) *PlainQuery {
	if q != nil && q.Command != "" {
		qslice := make([]Query, len(qs)+1)
		qslice[0] = q
		copy(qslice[1:], qs)
		*q = *(Or(qslice...))
	} else {
		*q = *(Or(qs...))
	}
	return q
}
func NewFromQuery() *FromQuery {
	return &FromQuery{
		Tables: [][2]string{},
	}

}

type FromQuery struct {
	Tables [][2]string
}

func (q *FromQuery) AddAlias(alias string, tableName string) *FromQuery {
	q.Tables = append(q.Tables, [2]string{tableName, alias})
	return q
}

func (q *FromQuery) Add(tableName string) *FromQuery {
	q.Tables = append(q.Tables, [2]string{tableName, ""})
	return q
}

func (q *FromQuery) QueryCommand() string {
	var command = ""
	if len(q.Tables) == 0 {
		return command
	}
	command = "FROM "
	for k := range q.Tables {
		if q.Tables[k][1] == "" {
			command += q.Tables[k][0] + " , "
		} else {
			command += q.Tables[k][0] + " as " + q.Tables[k][1] + " , "
		}
	}
	return command[:len(command)-3]
}
func (q *FromQuery) QueryArgs() []interface{} {
	return []interface{}{}
}

type JoinData struct {
	Type         string
	Table        [2]string
	Condition    *PlainQuery
	UsingColnums []string
}

func (d *JoinData) Alias(alias string, tableName string) *JoinData {
	d.Table[0] = tableName
	d.Table[1] = alias
	return d
}
func (d *JoinData) On(condition *PlainQuery) *JoinData {
	d.Condition = condition
	return d
}
func (d *JoinData) Using(colnums ...string) *JoinData {
	d.UsingColnums = colnums
	return d
}

func (d *JoinData) QueryCommand() string {
	var command = d.Type + " Join "
	command += d.Table[indexTableName]
	if d.Table[indexAlias] != "" {
		command += " AS " + d.Table[indexAlias]
	}
	if len(d.UsingColnums) == 0 {
		command += " ON " + d.Condition.QueryCommand()
	} else {
		command += " USING (" + strings.Join(d.UsingColnums, " , ") + ")"
	}
	return command
}
func (q *JoinData) QueryArgs() []interface{} {
	if q.Condition != nil && len(q.UsingColnums) == 0 {
		return q.Condition.QueryArgs()
	}
	return []interface{}{}
}

func NewJoinQuery() *JoinQuery {
	return &JoinQuery{
		Data: []*JoinData{},
	}
}

type JoinQuery struct {
	Data []*JoinData
}

func (q *JoinQuery) join(jointype string) *JoinData {
	data := &JoinData{
		Type:         jointype,
		Table:        [2]string{},
		Condition:    nil,
		UsingColnums: []string{},
	}
	q.Data = append(q.Data, data)
	return data
}

func (q *JoinQuery) InnerJoin() *JoinData {
	return q.join("INNER")
}
func (q *JoinQuery) LeftJoin() *JoinData {
	return q.join("LEFT")
}
func (q *JoinQuery) RightJoin() *JoinData {
	return q.join("RIGHT")
}

func (q *JoinQuery) QueryCommand() string {
	var command = ""
	for k := range q.Data {
		c := q.Data[k].QueryCommand()
		if c != "" {
			command += c + "\n"
		}
	}
	if command != "" {
		command = command[:len(command)-1]
	}
	return command
}
func (q *JoinQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		a := q.Data[k].QueryArgs()
		if a != nil {
			args = append(args, a...)
		}
	}
	return args
}

type QueryData struct {
	Field string
	Data  []interface{}
	Raw   string
}

func NewInsertQuery(tableName string) *InsertQuery {
	return &InsertQuery{
		Prefix:    New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type InsertQuery struct {
	Prefix    *PlainQuery
	TableName string
	Alias     string
	Data      []QueryData
}

func (q *InsertQuery) SetAlias(alias string) *InsertQuery {
	q.Alias = alias
	return q
}
func (q *InsertQuery) Add(field string, data interface{}) *InsertQuery {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		})
	return q
}
func (q *InsertQuery) AddRaw(field string, raw string) *InsertQuery {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}

func (q *InsertQuery) AddSelect(field string, Select *Select) *InsertQuery {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}

func (q *InsertQuery) QueryCommand() string {
	var command = "INSERT"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	if q.Alias != "" {
		command += " AS " + q.Alias
	}
	command += " INTO " + q.TableName
	var values = ""
	var columns = ""
	for k := range q.Data {
		if q.Data[k].Raw == "" {
			values += "? , "
		} else {
			values += q.Data[k].Raw + " , "
		}
		columns += q.Data[k].Field + " , "
	}
	command += " ("
	command += columns[:len(columns)-3]
	command += " )"

	command += " VALUES ( "
	command += values[:len(values)-3]
	command += " )"
	return command
}
func (q *InsertQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		if q.Data[k].Data != nil {
			args = append(args, q.Data[k].Data...)
		}
	}
	var result = []interface{}{}
	result = append(result, q.Prefix.QueryArgs()...)
	result = append(result, args...)
	return result
}

type DeleteQuery struct {
	TableName string
	Prefix    *PlainQuery
}

func (q *DeleteQuery) QueryCommand() string {
	var command = "DELETE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " FROM " + q.TableName
	return command
}

func (q *DeleteQuery) QueryArgs() []interface{} {
	return q.Prefix.QueryArgs()
}

func NewDeleteQuery(tableName string) *DeleteQuery {
	return &DeleteQuery{
		Prefix:    New(""),
		TableName: tableName,
	}
}

func NewUpdateQuery(tableName string) *UpdateQuery {
	return &UpdateQuery{
		Prefix:    New(""),
		TableName: tableName,
		Data:      []QueryData{},
	}
}

type UpdateQuery struct {
	Prefix    *PlainQuery
	TableName string
	Alias     string
	Data      []QueryData
}

func (q *UpdateQuery) SetAlias(alias string) *UpdateQuery {
	q.Alias = alias
	return q
}
func (q *UpdateQuery) AddSelect(field string, Select *Select) *UpdateQuery {
	query := *Select.Query()
	q.Data = append(q.Data, QueryData{
		Field: field,
		Raw:   "( " + query.QueryCommand() + " )",
		Data:  query.QueryArgs(),
	})
	return q
}
func (q *UpdateQuery) Add(field string, data interface{}) *UpdateQuery {
	q.Data = append(q.Data,
		QueryData{
			Field: field,
			Data:  []interface{}{data},
		},
	)
	return q
}
func (q *UpdateQuery) AddRaw(field string, raw string) *UpdateQuery {
	q.Data = append(q.Data, QueryData{Field: field, Raw: raw})
	return q
}
func (q *UpdateQuery) QueryCommand() string {
	var command = "UPDATE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " " + q.TableName
	if q.Alias != "" {
		command += " AS " + q.Alias
	}
	command += " SET "
	var values = ""
	for k := range q.Data {
		values += q.Data[k].Field + " = "
		if q.Data[k].Raw == "" {
			values += "? , "
		} else {
			values += q.Data[k].Raw + " , "
		}
	}
	command += values[:len(values)-3]
	return command
}
func (q *UpdateQuery) QueryArgs() []interface{} {
	var args = []interface{}{}
	for k := range q.Data {
		if q.Data[k].Data != nil {
			args = append(args, q.Data[k].Data...)
		}
	}
	var result = []interface{}{}
	result = append(result, q.Prefix.QueryArgs()...)
	result = append(result, args...)
	return result
}
func NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Prefix: New(""),
		Fields: []string{},
	}
}

type SelectQuery struct {
	Prefix *PlainQuery
	Fields []string
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
func (r *SelectResult) Args() []interface{} {
	return r.args
}

func (r *SelectResult) ScanFrom(s ResultScanner) error {
	return s.Scan(r.Args()...)
}
func NewWhereQuery() *WhereQurey {
	return &WhereQurey{
		Condition: New(""),
	}
}

type WhereQurey struct {
	Condition *PlainQuery
}

func (q *WhereQurey) QueryCommand() string {
	var command = q.Condition.QueryCommand()
	if command != "" {
		command = "WHERE " + command
	}
	return command
}
func (q *WhereQurey) QueryArgs() []interface{} {
	return q.Condition.QueryArgs()
}
func NewSelect() *Select {
	return &Select{
		Select: NewSelectQuery(),
		From:   NewFromQuery(),
		Join:   NewJoinQuery(),
		Where:  NewWhereQuery(),
		Other:  New(""),
	}
}

type Select struct {
	Select *SelectQuery
	From   *FromQuery
	Join   *JoinQuery
	Where  *WhereQurey
	Other  *PlainQuery
}

func (s *Select) Result() *SelectResult {
	return s.Select.Result()
}

func (s *Select) Query() *PlainQuery {
	return Lines(s.Select, s.From, s.Join, s.Where, s.Other)
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
func NewDelete(TableName string) *Delete {
	return &Delete{
		Delete: NewDeleteQuery(TableName),
		Where:  NewWhereQuery(),
		Other:  New(""),
	}
}

type Delete struct {
	Delete *DeleteQuery
	Where  *WhereQurey
	Other  *PlainQuery
}

func (d *Delete) Query() *PlainQuery {
	return Lines(d.Delete, d.Where, d.Other)
}

func NewInsert(tableName string) *Insert {
	return &Insert{
		Insert: NewInsertQuery(tableName),
		Other:  New(""),
	}
}

type Insert struct {
	Insert *InsertQuery
	Other  *PlainQuery
}

func (i *Insert) Query() *PlainQuery {
	return Lines(i.Insert, i.Other)
}
func NewUpdate(tableName string) *Update {
	return &Update{
		Update: NewUpdateQuery(tableName),
		Where:  NewWhereQuery(),
		Other:  New(""),
	}
}

type Update struct {
	Update *UpdateQuery
	Where  *WhereQurey
	Other  *PlainQuery
}

func (u *Update) Query() *PlainQuery {
	return Lines(u.Update, u.Where, u.Other)
}
