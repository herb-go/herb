package db

import (
	"database/sql"
)

type DB interface {
	SetDB(db *sql.DB)
	DB() *sql.DB
	TableName(string) string
}

type Table interface {
	DB() *sql.DB
	DBTableName() string
}

func New() *PlainDB {
	return &PlainDB{}
}

type PlainDB struct {
	db     *sql.DB
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

func (d *PlainDB) SetPrefix(prefix string) {
	d.prefix = prefix
}

func (d *PlainDB) Prefix() string {
	return d.prefix
}

func (d *PlainDB) TableName(tableName string) string {
	return d.prefix + tableName
}
func (d *PlainDB) Table(tableName string) *PlainTable {
	return NewTable(d, tableName)
}

func NewTable(db DB, tableName string) *PlainTable {
	return &PlainTable{
		db:    db,
		table: tableName,
	}
}

type PlainTable struct {
	db    DB
	table string
}

func (t *PlainTable) DB() *sql.DB {
	return t.db.DB()
}

func (t *PlainTable) SetTableName(table string) {
	t.table = table
}

func (t *PlainTable) TableName() string {
	return t.table
}

func (t *PlainTable) DBTableName() string {
	return t.db.TableName(t.table)
}
