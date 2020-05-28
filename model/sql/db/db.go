package db

import (
	"database/sql"
	"errors"
)

//ErrSetDriverFromTable error raised when execute SetDriver method of table .
var ErrSetDriverFromTable = errors.New("herb:sql/db you can't execute set driver method in table interface")

//Database database interface
type Database interface {
	//SetDB set sqlDB to database interface
	SetDB(db *sql.DB)
	//DB get sql DB of database.
	DB() *sql.DB
	//Driver return database drvier name.
	Driver() string
	//SetDriver set driver name.
	SetDriver(string)
	//BuildTableName return table name with giver table.
	BuildTableName(table string) string
	//Exec exec query with args.
	Exec(query string, args ...interface{}) (sql.Result, error)
	//Query exec query with args .
	//Return rows.
	Query(query string, args ...interface{}) (*sql.Rows, error)
	//QueryRow exec query with args and rows.
	//Return row.
	QueryRow(query string, args ...interface{}) *sql.Row
}

//Table table interface
type Table interface {
	Database
	//BuildFieldName return field name with given field.
	BuildFieldName(field string) string
	//SetAlias set table alias
	SetAlias(string)
	//Alias return table alias
	Alias() string
	//TableName return table name
	TableName() string
}

//New create new plain database.
func New() *PlainDB {
	return &PlainDB{}
}

//PlainDB plain database struct.
type PlainDB struct {
	db               *sql.DB
	driver           string
	prefix           string
	OptimizerFactory OptimizerFactory
	Optimizer        Optimizer
}

//Copy copy src plain db to dsc plain db
func Copy(src *PlainDB, dsc *PlainDB) {
	dsc.db = src.db
	dsc.driver = src.driver
	dsc.prefix = src.prefix
}

//Init init plain database with given option.
func (d *PlainDB) Init(o PlainDBOption) error {
	return o.ApplyTo(d)
}

//SetDB set sql db to plain datatbase.
func (d *PlainDB) SetDB(db *sql.DB) {
	d.db = db
}

//DB return sql db.
func (d *PlainDB) DB() *sql.DB {
	return d.db
}

//SetDriver set driver name.
func (d *PlainDB) SetDriver(driver string) {
	d.driver = driver
}

//Driver return driver name.
func (d *PlainDB) Driver() string {
	return d.driver
}

//SetPrefix set table prefix.
func (d *PlainDB) SetPrefix(prefix string) {
	d.prefix = prefix
}

//Prefix return table prefix.
func (d *PlainDB) Prefix() string {
	return d.prefix
}

//Exec exec query with args.
func (d *PlainDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if d.Optimizer != nil {
		query, args = d.Optimizer.MustOptimize(query, args)
	}
	return d.db.Exec(query, args...)
}

//Query exec query with args .
//Return rows.
func (d *PlainDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if d.Optimizer != nil {
		query, args = d.Optimizer.MustOptimize(query, args)
	}
	return d.db.Query(query, args...)
}

//QueryRow exec query with args and rows.
//Return row.
func (d *PlainDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if d.Optimizer != nil {
		query, args = d.Optimizer.MustOptimize(query, args)
	}
	return d.db.QueryRow(query, args...)
}

//BuildTableName return table name with giver table.
func (d *PlainDB) BuildTableName(tableName string) string {
	return d.prefix + tableName
}

//Table create plain table with given table name.
func (d *PlainDB) Table(tableName string) *PlainTable {
	return NewTable(d, tableName)
}

//NewTable create plain table with given database and table name.
func NewTable(db Database, tableName string) *PlainTable {
	return &PlainTable{
		Database: db,
		table:    tableName,
	}
}

//PlainTable plain table struct
type PlainTable struct {
	Database
	alias string
	table string
}

//SetName set plain table name.
func (t *PlainTable) SetName(table string) {
	t.table = table
}

//Name return plain table name.
func (t *PlainTable) Name() string {
	return t.table
}

//TableName return table name build with database.
func (t *PlainTable) TableName() string {
	return t.Database.BuildTableName(t.table)
}

//SetDriver painc if execute SetDriver method of  table.
func (t *PlainTable) SetDriver(driver string) {
	panic(ErrSetDriverFromTable)
}

//BuildFieldName build field name with alias.
func (t *PlainTable) BuildFieldName(name string) string {
	return t.Alias() + "." + name
}

//SetAlias set table alias.
func (t *PlainTable) SetAlias(alias string) {
	t.alias = alias
}

//Alias return table alias.
func (t *PlainTable) Alias() string {
	return t.alias
}
