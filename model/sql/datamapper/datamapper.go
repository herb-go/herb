package datamapper

import (
	"database/sql"
)

type DB interface {
	SetDB(db *sql.DB)
	DB() *sql.DB
	TableName(string) string
}

type DataMapper interface {
	DB() *sql.DB
	DBTableName() string
}

func NewDB(db *sql.DB, prefix string) *PrefixDB {
	return &PrefixDB{
		db:     db,
		prefix: prefix,
	}
}

type PrefixDB struct {
	db     *sql.DB
	prefix string
}

func (d *PrefixDB) SetDB(db *sql.DB) {
	d.db = db
}

func (d *PrefixDB) DB() *sql.DB {
	return d.db
}

func (d *PrefixDB) SetPrefix(prefix string) {
	d.prefix = prefix
}

func (d *PrefixDB) Prefix() string {
	return d.prefix
}

func (d *PrefixDB) TableName(tableName string) string {
	return d.prefix + tableName
}
func New(db DB, table string) *TableDataMapper {
	return &TableDataMapper{
		db:    db,
		table: table,
	}
}

type TableDataMapper struct {
	db    DB
	table string
}

func (t *TableDataMapper) DB() *sql.DB {
	return t.db.DB()
}

func (t *TableDataMapper) SetTableName(table string) {
	t.table = table
}

func (t *TableDataMapper) TableName() string {
	return t.table
}

func (t *TableDataMapper) DBTableName() string {
	return t.db.TableName(t.table)
}
