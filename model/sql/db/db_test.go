package db

import (
	"database/sql"
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const tablename = "testdb"

func TestDB(t *testing.T) {
	config := &Config{}
	err := json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	db := New()
	err = db.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.DB().Close()

	if db.Prefix() != config.Prefix {
		t.Error(db.Prefix())
	}
	if db.Driver() != db.Driver() {
		t.Error(db.Driver())
	}
	db.SetDriver("mysql")
	if db.Driver() != "mysql" {
		t.Error(db.Driver())
	}
	dm := db.Table("test")
	if dm.DB() != db.DB() {
		t.Error("db not equal")
	}
	if dm.Name() != "test" {
		t.Error(dm.Name())
	}
	if dm.TableName() != db.Prefix()+dm.Name() {
		t.Error(dm.TableName())
	}
	dm.SetName("testnew")
	if dm.Name() != "testnew" {
		t.Error(dm.TableName())
	}
	if dm.BuildFieldName("test") != dm.Alias()+".test" {
		t.Error(dm.BuildFieldName("test"))
	}
	dm.SetAlias("testdb")
	if dm.Alias() != "testdb" {
		t.Error(dm.Alias())
	}
	if dm.BuildFieldName("test") != "testdb.test" {
		t.Error(dm.BuildFieldName("test"))
	}
	defer func() {
		e := recover()
		if e == nil {
			t.Error(e)
		}
		err, ok := e.(error)
		if ok == false {
			t.Error(ok)
		}
		if err == nil {
			t.Error(err)
		}
	}()
	dm.SetDriver("sqlite")
	t.Error("Not panic")
}
func TestQuery(t *testing.T) {
	config := &Config{}
	err := json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	db := New()
	err = db.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.DB().Close()
	dm := db.Table(tablename)
	db.Exec("drop table " + dm.TableName() + ";")
	_, err = db.Exec("create table " + dm.TableName() + "( id varchar(255) );")
	if err != nil {
		t.Fatal(err)
	}
	_ = dm.QueryRow("select * from " + dm.TableName())
	if err != nil {
		t.Fatal(err)
	}
	rows, err := dm.Query("select * from " + dm.TableName())
	if err != nil {
		t.Fatal(err)
	}
	err = rows.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTxQuery(t *testing.T) {
	config := &Config{}
	err := json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	db := New()
	err = db.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.DB().Close()
	dt, err := NewTxDB(db)
	if err != nil {
		t.Fatal(err)
	}
	dm := NewTable(dt, tablename)
	db.Exec("drop table " + dm.TableName() + ";")
	_, err = dt.Exec("create table " + dm.TableName() + "( id varchar(255) );")
	if err != nil {
		t.Fatal(err)
	}
	err = dt.Commit()
	if err != nil {
		t.Fatal(err)
	}
	dt, err = NewTxDB(db)
	if err != nil {
		t.Fatal(err)
	}
	dm = NewTable(dt, tablename)
	row := dm.QueryRow("select * from " + dm.TableName())
	err = row.Scan()
	if err != sql.ErrNoRows {
		t.Fatal(err)
	}
	rows, err := dm.Query("select * from " + dm.TableName())
	if err != nil {
		t.Fatal(err)
	}
	err = rows.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = dt.Commit()
	if err != nil {
		t.Fatal(err)
	}
	dt, err = NewTxDB(db)
	if err != nil {
		t.Fatal(err)
	}
	err = dt.Rollback()
	if err != nil {
		t.Fatal(err)
	}
}
