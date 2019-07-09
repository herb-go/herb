package sqlitecolumn

import (
	"errors"
	"strings"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herbgo/column"
)

type Column struct {
	CID     int64
	Name    string
	Field   string
	Type    string
	NotNull string
	Default interface{}
	Key     string
}

func ConvertType(t string) (string, error) {
	ft := strings.Split(t, "(")[0]
	switch strings.ToUpper(ft) {
	case "BOOL":
		return "byte", nil
	case "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "TINYINT", "INT2", "INT8":
		return "int", nil
	case "BIGINT", "INT64":
		return "int64", nil
	case "FLOAT":
		return "float32", nil
	case "DOUBLE", "DOUBLE PRECISION", "REAL":
		return "float64", nil
	case "DATETIME", "DATE":
		return "time.Time", nil
	case "CHAR", "VARCHAR", "CHARACTER", "NCHAR", "NVARCHAR", "TEXT":
		return "string", nil
	case "BLOB":
		return "[]byte", nil
	}
	return "", errors.New("MysqlColumn:Column type " + t + " is not supported.")

}

func (c *Column) Convert() (*column.Column, error) {
	output := &column.Column{}
	output.Field = c.Field
	t, err := ConvertType(c.Type)
	output.ColumnType = t
	if err != nil {
		return nil, err
	}
	if output.ColumnType == "time.Time" && c.Default != nil {
		output.AutoValue = true
	}
	if (output.ColumnType == "int" || output.ColumnType == "int64") && c.Key == "1" {
		output.AutoValue = true
	}
	if c.Key == "1" {
		output.PrimayKey = true
	}
	if c.NotNull == "1" {
		output.NotNull = true
	}

	return output, nil
}

type Columns []Column

func (c *Columns) Columns() ([]column.Column, error) {
	output := []column.Column{}
	for _, v := range *c {
		column, err := v.Convert()
		if err != nil {
			return nil, err
		}
		output = append(output, *column)
	}

	return output, nil
}
func (c *Columns) Load(conn db.Database, table string) error {
	db := conn.DB()
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return err
	}
	defer rows.Close()
	*c = []Column{}
	for rows.Next() {
		column := Column{}
		if err := rows.Scan(&column.CID, &column.Field, &column.Type, &column.NotNull, &column.Default, &column.Key); err != nil {
			return err
		}
		*c = append(*c, column)
	}
	return nil
}

func init() {
	column.Drivers["sqlite3"] = func() column.ColumnsLoader {
		return &Columns{}
	}
}
