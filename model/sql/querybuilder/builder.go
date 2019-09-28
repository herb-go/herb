package querybuilder

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

//Builder query builder struct
type Builder struct {
	Driver string
	driver BuilderDriver
	lock   sync.Mutex
}

//IsDuplicate check if error is Is duplicate error.
func (b *Builder) IsDuplicate(err error) bool {
	return b.LoadDriver().IsDuplicate(err)
}

//TruncateTableQuery return truncate table query
func (b *Builder) TruncateTableQuery(table string) *PlainQuery {
	return b.New(b.LoadDriver().TruncateTableCommandBuilder(table))
}

// Exec exec query in given db.
//return sql result and error
//query will be convert by builder driver.
func (b *Builder) Exec(db DB, q Query) (sql.Result, error) {
	cmd, args := b.LoadDriver().ConvertQuery(q)
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

//QueryRow query rowsfrom db with given query.
//query will be convert by builder driver.
func (b *Builder) QueryRow(db DB, q Query) *sql.Row {
	cmd, args := b.LoadDriver().ConvertQuery(q)
	var timestamp int64
	if Debug {
		timestamp = time.Now().UnixNano()
	}
	row := db.QueryRow(cmd, args...)
	if Debug {
		Logger(timestamp, cmd, args)
	}
	return row
}

//QueryRows query rows from db with given query.
//query will be convert by builder driver.
func (b *Builder) QueryRows(db DB, q Query) (*sql.Rows, error) {
	cmd, args := b.LoadDriver().ConvertQuery(q)
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

//CountField return select count  field
func (b *Builder) CountField() string {
	return b.LoadDriver().CountField()
}

//New create new query builder
func New() *Builder {
	return &Builder{}
}

//LoadDriver load builder driver by Driver field
//Only load one time.
func (b *Builder) LoadDriver() BuilderDriver {
	d := b.driver
	if d != nil {
		return d
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	d = loadBuilderDriver(b.Driver)
	b.driver = d
	return b.driver
}

//DefaultDriver default driver
var DefaultDriver = &EmptyBuilderDriver{}

//BuilderDriver query builder driver interface
type BuilderDriver interface {
	ConvertQuery(q Query) (string, []interface{})
	TruncateTableCommandBuilder(t string) string
	LimitCommandBuilder(q *LimitClause) string
	LimitArgBuilder(q *LimitClause) []interface{}
	DeleteCommandBuilder(q *DeleteClause) string
	DeleteArgBuilder(q *DeleteClause) []interface{}
	CountField() string
	IsDuplicate(error) bool
}

// EmptyBuilderDriver empty query builder.
// Using mysql statements
type EmptyBuilderDriver struct {
}

//TruncateTableCommandBuilder return truncate table query.
func (d EmptyBuilderDriver) TruncateTableCommandBuilder(t string) string {
	return "TRUNCATE TABLE " + t
}

//ConvertQuery  convert query to command and args
func (d EmptyBuilderDriver) ConvertQuery(q Query) (string, []interface{}) {
	return q.QueryCommand(), q.QueryArgs()
}

//CountField return count field
func (d *EmptyBuilderDriver) CountField() string {
	return "count(*)"
}

//LimitCommandBuilder build limit command with given limit query.
func (d *EmptyBuilderDriver) LimitCommandBuilder(q *LimitClause) string {
	var command = ""
	if q.Limit != nil {
		command = "LIMIT ? "
	}
	if q.Offset != nil {
		command += "OFFSET ? "
	}
	return command
}

//LimitArgBuilder build limit args with given limit clause.
func (d *EmptyBuilderDriver) LimitArgBuilder(q *LimitClause) []interface{} {
	var args = []interface{}{}
	if q.Limit != nil {
		args = append(args, *q.Limit)
	}
	if q.Offset != nil {
		args = append(args, *q.Offset)
	}
	return args
}

// DeleteCommandBuilder build delete command  with given delete clause.
func (d *EmptyBuilderDriver) DeleteCommandBuilder(q *DeleteClause) string {
	var command = "DELETE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	command += " FROM " + q.TableName
	return command
}

// DeleteArgBuilder build delete args  with given delete clause.
func (d *EmptyBuilderDriver) DeleteArgBuilder(q *DeleteClause) []interface{} {
	return q.Prefix.QueryArgs()
}

//IsDuplicate check if error is Is duplicate error.
func (d *EmptyBuilderDriver) IsDuplicate(err error) bool {
	return false
}

var drivers = map[string]BuilderDriver{}

//RegisterDriver register querybuilder driver with given  name.
func RegisterDriver(name string, driver BuilderDriver) {
	drivers[name] = driver
	if Debug && driver != nil {
		fmt.Println("querybuilder: build driver '" + name + "' registered.")
	}
}

func loadBuilderDriver(name string) BuilderDriver {
	d := drivers[name]
	if d == nil {
		return DefaultDriver
	}
	return d
}
