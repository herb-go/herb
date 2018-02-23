package datamapper

import "database/sql"

type PrefixDBOption interface {
	Apply(*PrefixDB) error
}
type DBConfig struct {
	Driver string
	Conn   string
	Prefix string
}

func (c *DBConfig) Apply(d *PrefixDB) error {
	db, err := sql.Open(c.Driver, c.Conn)
	if err != nil {
		return err
	}
	d.SetDB(db)
	d.SetPrefix(c.Prefix)
	return nil
}
