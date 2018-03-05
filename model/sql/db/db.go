package db

import (
	"database/sql"
)

type Database interface {
	SetDB(db *sql.DB)
	DB() *sql.DB
	Driver() string
	SetDriver(string)
	BuildTableName(string) string
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Table interface {
	Database
	BuildFieldName(name string) string
	SetAlias(string)
	Alias() string
	TableName() string
}

func New() *PlainDB {
	return &PlainDB{}
}

type PlainDB struct {
	db     *sql.DB
	driver string
	prefix string
}

func (d *PlainDB) Init(o PlainDBOption) error {
	return o.Apply(d)
}
func (d *PlainDB) SetDB(db *sql.DB) {
	d.db = db
}

func (d *PlainDB) DB() *sql.DB {
	return d.db
}
func (d *PlainDB) SetDriver(driver string) {
	d.driver = driver
}

func (d *PlainDB) Driver() string {
	return d.driver
}
func (d *PlainDB) SetPrefix(prefix string) {
	d.prefix = prefix
}

func (d *PlainDB) Prefix() string {
	return d.prefix
}

func (d *PlainDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}
func (d *PlainDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}
func (d *PlainDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

func (d *PlainDB) BuildTableName(tableName string) string {
	return d.prefix + tableName
}
func (d *PlainDB) Table(tableName string) *PlainTable {
	return NewTable(d, tableName)
}

func NewTable(db Database, tableName string) *PlainTable {
	return &PlainTable{
		Database: db,
		table:    tableName,
	}
}

type PlainTable struct {
	Database
	alias string
	table string
}

func (t *PlainTable) DB() *sql.DB {
	return t.Database.DB()
}

func (t *PlainTable) SetName(table string) {
	t.table = table
}

func (t *PlainTable) Name() string {
	return t.table
}

func (t *PlainTable) TableName() string {
	return t.Database.BuildTableName(t.table)
}

func (t *PlainTable) BuildFieldName(name string) string {
	if t.alias == "" {
		return name
	}
	return t.alias + "." + name
}
func (t *PlainTable) SetAlias(alias string) {
	t.alias = alias
}
func (t *PlainTable) Alias() string {
	return t.alias
}
