package db

import (
	"database/sql"
	"time"
)

//DefaultConnMaxLifetimeInSecond default conn max lifetime
const DefaultConnMaxLifetimeInSecond = int64(30)

//PlainDBOption plain database init option interface.
type PlainDBOption interface {
	//Apply init plain database.
	Apply(*PlainDB) error
}

//Config database config
type Config struct {
	//Driver sql driver.
	Driver string
	//Type sql database type.
	Type string
	//Conn sql conn string.
	DataSource string
	//Prefix sql table prefix.
	Prefix string
	//MaxIdleConns max idle conns.
	MaxIdleConns int
	//ConnMaxLifetimeInSecond conn max Lifetime in second.
	ConnMaxLifetimeInSecond int64
	//MaxOpenConns max open conns.
	MaxOpenConns int
}

//Apply init plain database with config
func (c *Config) Apply(d *PlainDB) error {
	db, err := sql.Open(c.Driver, c.DataSource)
	if err != nil {
		return err
	}
	if c.MaxIdleConns > 0 {
		db.SetMaxIdleConns(c.MaxIdleConns)
	}
	if c.ConnMaxLifetimeInSecond > 0 {
		db.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetimeInSecond) * time.Second)
	} else if c.ConnMaxLifetimeInSecond == 0 {
		db.SetConnMaxLifetime(time.Duration(DefaultConnMaxLifetimeInSecond) * time.Second)
	}
	if c.MaxOpenConns > 0 {
		db.SetMaxOpenConns(c.MaxOpenConns)
	}
	d.SetDB(db)
	d.SetPrefix(c.Prefix)
	return nil
}
