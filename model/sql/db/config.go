package db

import "database/sql"

type PlainDBOption interface {
	Apply(*PlainDB) error
}
type DBConfig struct {
	Driver string
	Conn   string
	Prefix string
}

func (c *DBConfig) Apply(d *PlainDB) error {
	db, err := sql.Open(c.Driver, c.Conn)
	if err != nil {
		return err
	}
	d.SetDB(db)
	d.SetPrefix(c.Prefix)
	return nil
}
