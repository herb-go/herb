package posgresqlcolumn

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq" //mysql driver

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/db/columns"
)

//Column mysql column struct
type Column struct {
	// Field  column name
	Field string
	// Type column type
	Type string
	// IsNull if column can be null
	IsNull string
	// Key  if column is primary key
	Key bool
	// Default default value
	Default interface{}
	// Extra extra data
	Extra string
}

// ConvertType convert culumn type to golang type.
func ConvertType(t string) (string, error) {
	ft := strings.Split(t, "(")[0]
	switch strings.ToUpper(ft) {

	case "BOOLEAN":
		return "bool", nil
	case "SMALLINT", "INT", "INTEGER":
		return "int", nil
	case "BIGINT":
		return "int64", nil
	case "REAL":
		return "float32", nil
	case "DOUBLE PRECISION":
		return "float64", nil
	case "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITHOUT TIME ZONE":
		return "time.Time", nil
	case "CHARACTER", "CHARACTER VARYING", "TEXT", "JSON", "XML", "UUID":
		return "string", nil
	case "BYTEA", "JSONB":
		return "[]byte", nil
	}
	return "", errors.New("postgresqlColumn:column type " + t + " is not supported.")

}

// Convert convert MysqlColumn to commn column
func (c *Column) Convert() (*columns.Column, error) {
	output := &columns.Column{}
	output.Field = c.Field
	t, err := ConvertType(c.Type)
	output.ColumnType = t
	if err != nil {
		return nil, err
	}
	if c.Default != nil {
		output.AutoValue = true
	}
	if c.Key {
		output.PrimayKey = true
	}
	if c.IsNull == "NO" {
		output.NotNull = true
	}

	return output, nil
}

// Columns mysql columns type
type Columns []Column

// Columns return loaded columns
func (c *Columns) Columns() ([]*columns.Column, error) {
	output := []*columns.Column{}
	for _, v := range *c {
		column, err := v.Convert()
		if err != nil {
			return nil, err
		}
		output = append(output, column)
	}
	return output, nil
}

// Load load columns with given database and table name
func (c *Columns) Load(conn db.Database, table string) error {
	db := conn.DB()
	rows, err := db.Query(fmt.Sprintf("SELECT column_name ,data_type , is_nullable , column_default FROM information_schema.columns WHERE table_name = '%s' and table_schema='public';",
		table))
	if err != nil {
		return err
	}
	defer rows.Close()
	*c = []Column{}
	for rows.Next() {
		column := Column{}
		if err := rows.Scan(&column.Field, &column.Type, &column.IsNull, &column.Default); err != nil {
			return err
		}
		*c = append(*c, column)
	}
	rows, err = db.Query(fmt.Sprintf("SELECT column_name FROM information_schema.key_column_usage WHERE table_name = '%s' ",
		table))
	if err != nil {
		return err
	}
	pks := []string{}
	for rows.Next() {
		field := ""
		if err := rows.Scan(&field); err != nil {
			return err
		}
		pks = append(pks, field)
	}
	for _, v := range pks {
		for k := range *c {
			if (*c)[k].Field == v {
				(*c)[k].Key = true
				break
			}
		}
	}
	return nil
}

func init() {
	columns.Register("postgres", func() columns.Loader {
		return &Columns{}
	})
}
