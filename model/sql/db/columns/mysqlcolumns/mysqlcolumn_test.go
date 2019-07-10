package mysqlcolumns

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/db/columns"
)

var target = []*columns.Column{
	&columns.Column{Field: "id", ColumnType: "int", AutoValue: true, PrimayKey: true, NotNull: true},
	&columns.Column{Field: "f_nullable", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: false},
	&columns.Column{Field: "f_tinyint", ColumnType: "byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_bit", ColumnType: "byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_bool", ColumnType: "byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_smallint", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_mediumint", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_int", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_integer", ColumnType: "int", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_biginteger", ColumnType: "int64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_float", ColumnType: "float32", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_double", ColumnType: "float64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_doubleprecision", ColumnType: "float64", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_datetime", ColumnType: "time.Time", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_timestamp", ColumnType: "time.Time", AutoValue: true, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_char", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_varchar", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_tinytext", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_text", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_mediumtext", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_longtext", ColumnType: "string", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_binary", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_varbinary", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_tinyblob", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_blob", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_mediumblob", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
	&columns.Column{Field: "f_longblob", ColumnType: "[]byte", AutoValue: false, PrimayKey: false, NotNull: true},
}

func TestColumns(t *testing.T) {
	config := db.NewConfig()
	err := json.Unmarshal([]byte(ConfigJSON), config)
	if err != nil {
		t.Fatal(err)
	}
	db := db.New()
	err = db.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.DB().Close()
	c := columns.Drivers["mysql"]()
	err = c.Load(db, "columns")
	if err != nil {
		t.Fatal(err)
	}
	columns, err := c.Columns()
	if err != nil {
		t.Fatal(err)
	}

	if len(columns) != len(target) {
		t.Fatal(columns)
	}

	for k := range columns {
		if columns[k].Field != target[k].Field ||
			columns[k].ColumnType != target[k].ColumnType ||
			columns[k].AutoValue != target[k].AutoValue ||
			columns[k].PrimayKey != target[k].PrimayKey ||
			columns[k].NotNull != target[k].NotNull {
			t.Fatal(columns[k])
		}
	}
}

func TestConvertType(t *testing.T) {
	_, err := ConvertType("unknowtype")
	if err == nil || !strings.Contains(err.Error(), "unknowtype") {
		t.Fatal(err)
	}
}
