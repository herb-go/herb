package db

import (
	"database/sql"
)

// NewTxDB create new database wtih transaction by given database.
func NewTxDB(database Database) (*TxDB, error) {
	db := database.DB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &TxDB{Database: database, Tx: tx}, nil
}

//TxDB database wtih transaction.
type TxDB struct {
	Database
	Tx *sql.Tx
}

//Exec exec query with args.
func (d *TxDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.Tx.Exec(query, args...)
}

//Query exec query with args .
//Return rows.
func (d *TxDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.Tx.Query(query, args...)
}

//QueryRow exec query with args and rows.
//Return row.
func (d *TxDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.Tx.QueryRow(query, args...)
}

// Commit commits the transaction.
func (d *TxDB) Commit() error {
	return d.Tx.Commit()
}

// Rollback aborts the transaction.
func (d *TxDB) Rollback() error {
	return d.Tx.Rollback()
}
