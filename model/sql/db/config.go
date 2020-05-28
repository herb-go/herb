package db

import (
	"database/sql"
	"time"
)

//DefaultConnMaxLifetimeInSecond default conn max lifetime
const DefaultConnMaxLifetimeInSecond = int64(30)

//PlainDBOption plain database init option interface.
type PlainDBOption interface {
	//ApplyTo init plain database.
	ApplyTo(*PlainDB) error
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
	Optimizer    func(v interface{}) error `config:", lazyload"`
}

//ApplyTo init plain database with config
func (c *Config) ApplyTo(d *PlainDB) error {
	factories := Factories()
	for _, f := range factories {
		if f == c.Driver {
			driver, err := NewDriver(c.Driver, c)
			if err != nil {
				return err
			}
			return driver.ApplyTo(d)
		}
	}
	f := d.OptimizerFactory
	if f == nil {
		f = DefaultOptimizerFactory
	}
	o, err := f(c.Optimizer)
	if err != nil {
		return err
	}
	d.Optimizer = o
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
	if c.Type != "" {
		d.SetDriver(c.Type)
	} else {
		d.SetDriver(c.Driver)
	}
	return nil
}

//NewConfig create new config
func NewConfig() *Config {
	return &Config{}
}
