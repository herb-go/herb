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

//NewBuilder create new query builder
func NewBuilder() *Builder {
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
	LimitCommandBuilder(q *LimitQuery) string
	LimitArgBuilder(q *LimitQuery) []interface{}
	DeleteCommandBuilder(q *DeleteQuery) string
	DeleteArgBuilder(q *DeleteQuery) []interface{}
	CountField() string
}

// EmptyBuilderDriver empty query builder.
// Using mysql statements
type EmptyBuilderDriver struct {
}

func (d EmptyBuilderDriver) ConvertQuery(q Query) (string, []interface{}) {
	return q.QueryCommand(), q.QueryArgs()
}

//CountField return count field
func (d *EmptyBuilderDriver) CountField() string {
	return "count(*)"
}

//LimitCommandBuilder build limit command with given limit query.
func (d *EmptyBuilderDriver) LimitCommandBuilder(q *LimitQuery) string {
	var command = ""
	if q.limit != nil {
		command = "LIMIT ? "
	}
	if q.offset != nil {
		command += "OFFSET ? "
	}
	return command
}

//LimitArgBuilder build limit args with given limit query.
func (d *EmptyBuilderDriver) LimitArgBuilder(q *LimitQuery) []interface{} {
	var args = []interface{}{}
	if q.limit != nil {
		args = append(args, *q.limit)
	}
	if q.offset != nil {
		args = append(args, *q.offset)
	}
	return args
}

func (d *EmptyBuilderDriver) DeleteCommandBuilder(q *DeleteQuery) string {
	var command = "DELETE"
	p := q.Prefix.QueryCommand()
	if p != "" {
		command += " " + p
	}
	if q.alias != "" {
		command += " " + q.alias
	}
	command += " FROM " + q.TableName
	if q.alias != "" {
		command += " AS " + q.alias
	}
	return command
}
func (d *EmptyBuilderDriver) DeleteArgBuilder(q *DeleteQuery) []interface{} {
	return q.Prefix.QueryArgs()
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

func init() {
	RegisterDriver("mysql", DefaultDriver)
}
