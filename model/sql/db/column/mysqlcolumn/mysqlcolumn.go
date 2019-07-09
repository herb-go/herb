package mysqlcolumn

import (
	"errors"
	"strings"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herbgo/column"
)

type Column struct {
	Field   string
	Type    string
	IsNull  string
	Key     string
	Default interface{}
	Extra   string
}

func ConvertType(t string) (string, error) {
	ft := strings.Split(t, "(")[0]
	switch strings.ToUpper(ft) {
	case "TINYINT", "BIT", "BOOL":
		return "byte", nil
	case "SMALLINT", "MEDIUMINT", "INT", "INTEGER":
		return "int", nil
	case "BIGINT":
		return "int64", nil
	case "FLOAT":
		return "float32", nil
	case "DOUBLE", "DOUBLE PRECISION":
		return "float64", nil
	case "DATETIME", "TIMESTAMP":
		return "time.Time", nil
	case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT":
		return "string", nil
	case "BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
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
	if strings.Contains(c.Extra, "auto_increment") {
		output.AutoValue = true
	}
	if strings.Contains(c.Key, "PRI") {
		output.PrimayKey = true
	}
	if c.IsNull == "NO" {
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
	rows, err := db.Query("desc " + table)
	if err != nil {
		return err
	}
	defer rows.Close()
	*c = []Column{}
	for rows.Next() {
		column := Column{}
		if err := rows.Scan(&column.Field, &column.Type, &column.IsNull, &column.Key, &column.Default, &column.Extra); err != nil {
			return err
		}
		*c = append(*c, column)
	}
	return nil
}

func init() {
	column.Drivers["mysql"] = func() column.ColumnsLoader {
		return &Columns{}
	}
}
