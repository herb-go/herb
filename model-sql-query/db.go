package query

import "database/sql"

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func NewDB() *sql.DB {
	return &sql.DB{}
}

func LoadDB(db *sql.DB, driverName, dataSourceName string) error {
	opened, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	*db = *opened
	return nil
}

func MustLoadDB(db *sql.DB, driverName, dataSourceName string) {
	err := LoadDB(db, driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
}

type DBConfig struct {
	Driver string
	Conn   string
	Prefix string
}

func (d *DBConfig) Open() (*sql.DB, error) {
	return sql.Open(d.Driver, d.Conn)
}

func (d *DBConfig) MustOpen() *sql.DB {
	db, err := d.Open()
	if err != nil {
		panic(err)
	}
	return db
}

func (d *DBConfig) Load(db *sql.DB) error {
	return LoadDB(db, d.Driver, d.Conn)
}

func (d *DBConfig) MustLoad(db *sql.DB) {
	MustLoadDB(db, d.Driver, d.Conn)
}
